package contract_service

import (
	"io"
	"relocate/api"
	"relocate/model"
	"relocate/service/declaration_service"
	"relocate/service/staging_service"
	"relocate/util/convert"
	"relocate/util/errors"
	"relocate/util/excel"
	"strings"

	"github.com/jinzhu/gorm"
)

func UpdateHouseWriteOffStatusList(ContractNoList []string, HouseWriteOff bool, operator, data string) (err error) {
	err = model.UpdateHouseWriteOffStatusList(ContractNoList, HouseWriteOff)
	if err == nil {
		id, _ := model.GetNowStagingConfig()
		staging, _ := staging_service.GetStagingInfoById(id)
		logging := model.Logging{
			Username:    operator,
			StagingName: staging.StagingName,
			Operation:   "批量修改合同号的注销状态",
			Details:     data,
		}
		logging.Create()
	}
	return err
}

func GetStagingById(stagingID uint) bool {
	_, err := staging_service.GetStagingInfoById(stagingID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
	}
	return true
}

func AddContract(body api.AddContractBody, operator, data string) error {
	cardNumber := strings.Replace(body.CardNumber, "，", ",", -1)
	cardNumber = strings.Replace(cardNumber, "（", "(", -1)
	cardNumber = strings.Replace(cardNumber, "）", ")", -1)
	if body.Registration == "" {
		body.Registration = "无证"
	}
	if body.Proprietor == "" {
		body.Proprietor = "无证"
	}
	if body.CollectiveLandPropertyCertificate == "" {
		body.CollectiveLandPropertyCertificate = "无证"
	}
	contract := &model.Contract{
		ContractNo:                        body.ContractNo,
		SocialCategory:                    body.SocialCategory,
		Peoples:                           body.Peoples,
		CardNumber:                        cardNumber,
		HouseNumber:                       body.HouseNumber,
		OldAddress:                        body.OldAddress,
		PhoneNumber1:                      body.PhoneNumber1,
		PhoneNumber2:                      body.PhoneNumber2,
		DateOfSigning:                     body.DateOfSigning,
		DateOfDelivery:                    body.DateOfDelivery,
		Signatory:                         body.Signatory,
		Registration:                      body.Registration,
		InitialHQArea:                     body.InitialHQArea,
		RemainingHQArea:                   body.InitialHQArea,
		IsDelivery:                        body.IsDelivery,
		CanDeclare:                        body.CanDeclare,
		StagingID:                         body.StagingID,
		HouseWriteOff:                     body.HouseWriteOff,
		Proprietor:                        body.Proprietor,
		Desc:                              body.Desc,
		CollectiveLandPropertyCertificate: body.CollectiveLandPropertyCertificate,
		ChangeMethod:                      body.ChangeMethod,
	}
	contract.Single = staging_service.JudgmentSingle(cardNumber)
	if err := contract.Create(); err != nil {
		return err
	}
	id, _ := model.GetNowStagingConfig()
	staging, _ := staging_service.GetStagingInfoById(id)
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   "添加合同号：" + body.ContractNo,
		Details:     data,
	}
	logging.Create()
	return nil
}

func UpdateCardNumber(contractNo, cardNumber, operator, data string) error {
	if err := model.UpdateCardNumber(contractNo, cardNumber); err != nil {
		return err
	}
	id, _ := model.GetNowStagingConfig()
	staging, _ := staging_service.GetStagingInfoById(id)
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   "合同号：" + contractNo + "添加身份证",
		Details:     data,
	}
	logging.Create()
	return nil
}

func SupplementContract(r io.Reader) (int, error) {
	contracts, err := parseExcel(r)
	if err != nil {
		return 0, err
	}
	for _, contract := range *contracts {
		_ = model.UpdateTargetPlacementAreaAndTemporaryRelocationArea(contract.ContractNo, contract.TargetPlacementArea, contract.TemporaryRelocationArea)
	}
	return len(*contracts), nil
}

func parseExcel(r io.Reader) (*[]model.Contract, error) {
	var contracts []model.Contract
	exFile, err := excel.Open(r)
	if err != nil {
		return nil, err
	}
	results, err := exFile.GetSheetData("668355") //表名记号
	if err != nil {
		return nil, err
	}
	resultsLen := len(results)
	indexColumn := make(map[string]int)
	var initialLine int
	for initialLine = 0; initialLine < resultsLen; initialLine++ {
		for index, column := range results[initialLine] {
			//如果找到合同号关键字，则将此行记录识别为头部标题
			switch column {
			case "拆迁安置补偿协议合同号":
				indexColumn["ContractNo"] = index
				break
			case "指标安置面积":
				indexColumn["TargetPlacementArea"] = index
				break
			case "计算临迁费面积":
				indexColumn["TemporaryRelocationArea"] = index
				break
			}
		}
		if len(indexColumn) > 0 {
			initialLine++
			break
		}
	}
	if initialLine >= resultsLen {
		return nil, errors.BadError("无效数据")
	}
	for i := initialLine; i < resultsLen; i++ {
		if len(results[i]) > 0 {
			targetPlacementArea, _ := convert.StrToFloat64(results[i][indexColumn["TargetPlacementArea"]], 4)
			temporaryRelocationArea, _ := convert.StrToFloat64(results[i][indexColumn["TemporaryRelocationArea"]], 4)
			contracts = append(contracts, model.Contract{
				ContractNo:              results[i][indexColumn["ContractNo"]],
				TargetPlacementArea:     targetPlacementArea,
				TemporaryRelocationArea: temporaryRelocationArea,
			})
		}
	}
	return &contracts, nil
}

func getContractByStagingId(staginid uint) (*model.Contract, error) {
	var c model.Contract
	contract, err := c.FindContractByStagingId(staginid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("该合同单不存在")
		}
		return nil, err
	}
	return contract, nil
}

func SettingCanDeclareByStagingId(stagingid uint, status int) error {
	if _, err := getContractByStagingId(stagingid); err != nil {
		return err
	}
	if _, err := declaration_service.GetDeclarationByStaginId(stagingid); err != nil {
		return err
	}
	if err := model.UpdateCanDeclarationStatusByStaginId(stagingid, status); err != nil {
		return errors.BadError("修改申报状态失败")
	}
	return nil
}

func UpdateContract(canDeclare bool, changeMethod string,
	collectiveLandPropertyCertificate string, contractNo string, dateOfDelivery string,
	dateOfSigning string, desc string, houseNumber string, houseWriteOff bool, oldAddress string, phoneNumber1 string,
	phoneNumber2 string, proprietor string, registration string, signatory string, socialCategory string, stagingID uint) error {
	contract := model.Contract{
		CanDeclare:                        canDeclare,
		ChangeMethod:                      changeMethod,
		CollectiveLandPropertyCertificate: collectiveLandPropertyCertificate,
		DateOfDelivery:                    dateOfDelivery,
		DateOfSigning:                     dateOfSigning,
		Desc:                              desc,
		HouseNumber:                       houseNumber,
		HouseWriteOff:                     houseWriteOff,
		OldAddress:                        oldAddress,
		PhoneNumber1:                      phoneNumber1,
		PhoneNumber2:                      phoneNumber2,
		Proprietor:                        proprietor,
		Registration:                      registration,
		Signatory:                         signatory,
		SocialCategory:                    socialCategory,
		StagingID:                         stagingID,
	}
	if _, err := GetContractById(contractNo); err != nil {
		return err
	}
	if err := contract.UpdateContract(contractNo); err != nil {
		return errors.BadError("修改合同单失败")
	}

	return nil
}

func GetContractById(id string) (*model.Contract, error) {
	contract, err := model.FindContractById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("该合同单不存在")
		}
		return nil, err
	}
	return contract, nil
}

func GetInContractById(id []string) (*model.Contract, error) {
	contract, err := model.FindInContractById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("该合同单不存在")
		}
		return nil, err
	}
	return contract, nil
}

func UpdateCanDeclare(contractno []string, canDeclare bool) error {
	if _, err := GetInContractById(contractno); err != nil {
		return err
	}
	if err := model.UpdateCanDeclare(contractno, canDeclare); err != nil {
		return errors.BadError("（批量）修改状态失败")
	}
	return nil
}
