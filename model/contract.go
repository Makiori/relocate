package model

import (
	"bytes"
	"fmt"
	"relocate/util/logging"
	"relocate/util/times"
	"strconv"
)

// 回迁户合同单
type Contract struct {
	ContractNo     string `json:"contract_no" gorm:"primary_key;comment:'合同号-主键'"`
	SocialCategory string `json:"social_category" gorm:"null;comment:'社别'"`
	Peoples        string `json:"peoples" gorm:"not null;comment:'被拆迁人(可能有多人)'"`
	CardNumber     string `json:"card_number" gorm:"null;comment:'被拆迁人身份证号码(可能有多人)'"`
	HouseNumber    string `json:"house_number" gorm:"null;comment:'房屋栋号'"`
	OldAddress     string `json:"old_address" gorm:"null;comment:'被拆迁房屋地址'"`
	PhoneNumber1   string `json:"phone_number1" gorm:"null;comment:'手机号码1'"`
	PhoneNumber2   string `json:"phone_number2" gorm:"null;comment:'手机号码2'"`
	//Trustee                           string          `json:"trustee" gorm:"null;comment:'受托人'"`
	//TrusteeCardNumber                 string          `json:"trustee_card_number" gorm:"null;comment:'受托人身份证号码'"`
	//TrusteePhoneNumber                string          `json:"trustee_phone_number" gorm:"null;comment:'受托人手机号码'"`
	//TrusteeRelationship               string          `json:"trustee_relationship" gorm:"null;comment:'受托人关系'"`
	DateOfSigning                     string          `json:"date_of_signing" gorm:"null;comment:'签署协议日期'"`
	DateOfDelivery                    string          `json:"date_of_delivery" gorm:"null;comment:'交楼日期'"`
	Proprietor                        string          `json:"proprietor" gorm:"null;comment:'证载产权人'"`
	Signatory                         string          `json:"signatory" gorm:"null;comment:'签约人'"`
	ChangeMethod                      string          `json:"change_method" gorm:"null;comment:'变更方式'"`
	CollectiveLandPropertyCertificate string          `json:"collective_land_property_certificate" gorm:"null;comment:'集体土地房产证字'"`
	Registration                      string          `json:"registration" gorm:"null;comment:'登记字号'"`
	InitialHQArea                     float64         `json:"initial_hq_area" gorm:"null;comment:'初始回迁面积'"`
	RemainingHQArea                   float64         `json:"remaining_hq_area" gorm:"null;comment:'剩余回迁面积'"`
	Desc                              string          `json:"desc" gorm:"null;comment:'备注'"`
	IsDelivery                        bool            `json:"is_delivery" gorm:"null;comment:'是否交齐楼'"`
	HouseWriteOff                     bool            `json:"house_write_off" gorm:"null;comment:'是否完成房屋注销'"`
	CanDeclare                        bool            `json:"can_declare" gorm:"not null;comment:'可否申报'"`
	Single                            bool            `json:"single" gorm:"not null;comment:'单人、多人|单人直接where = |多人like匹配 身份证'"`
	StagingID                         uint            `json:"staging_id" gorm:"null;comment:'分期数ID'"`
	CreatedAt                         times.JsonTime  `json:"created_at" gorm:"not null;comment:'创建时间'"`
	UpdatedAt                         times.JsonTime  `json:"-" gorm:"comment:'更新时间'"`
	DeletedAt                         *times.JsonTime `json:"-" gorm:"comment:'删除时间'"`
	TargetPlacementArea               float64         `json:"target_placement_area" gorm:"null;comment:'指标安置面积'"`
	TemporaryRelocationArea           float64         `json:"temporary_relocation_area" gorm:"null;comment:'计算临迁费面积'"`
}

func (c Contract) TableName() string {
	return "contract"
}

func initContract() {
	if !db.HasTable(&Contract{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Contract{}).Error; err != nil {
			panic(err)
		}
	}
}

func (c *Contract) Create() error {
	return db.Create(&c).Error
}

func UpdateContractStagingID(contractNo string, stagingID uint) (err error) {
	return db.Model(&Contract{}).Where("contract_no = ?", contractNo).Update("staging_id", stagingID).Error
}

//当前版本gorm还不支持批量插入，手写sql进行批量插入
func BatchCreateContract(cs []Contract) error {
	var buffer bytes.Buffer
	sql := "insert into `contract` " +
		"(`contract_no`,`social_category`,`peoples`,`card_number`,`house_number`,`old_address`," +
		"`phone_number1`,`phone_number2`,`date_of_signing`,`date_of_delivery`,`proprietor`,`signatory`," +
		"`change_method`,`collective_land_property_certificate`,`registration`,`initial_hq_area`,`remaining_hq_area`," +
		"`desc`,`is_delivery`,`house_write_off`,`can_declare`,`single`,`staging_id`,`created_at`) " +
		"values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, c := range cs {
		if i == len(cs)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',%v,%v,%v,%v,%d,'%s');",
				c.ContractNo, c.SocialCategory, c.Peoples, c.CardNumber, c.HouseNumber, c.OldAddress, c.PhoneNumber1,
				c.PhoneNumber2, c.DateOfSigning, c.DateOfDelivery, c.Proprietor, c.Signatory, c.ChangeMethod, c.CollectiveLandPropertyCertificate,
				c.Registration, c.InitialHQArea, c.RemainingHQArea, c.Desc, c.IsDelivery, c.HouseWriteOff, c.CanDeclare,
				c.Single, c.StagingID, times.ToStr()))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',%v,%v,%v,%v,%d,'%s'),",
				c.ContractNo, c.SocialCategory, c.Peoples, c.CardNumber, c.HouseNumber, c.OldAddress, c.PhoneNumber1,
				c.PhoneNumber2, c.DateOfSigning, c.DateOfDelivery, c.Proprietor, c.Signatory, c.ChangeMethod, c.CollectiveLandPropertyCertificate,
				c.Registration, c.InitialHQArea, c.RemainingHQArea, c.Desc, c.IsDelivery, c.HouseWriteOff, c.CanDeclare,
				c.Single, c.StagingID, times.ToStr()))
		}
	}
	return db.Exec(buffer.String()).Error
}

type ContractLikeBody struct {
	ContractNo                        string  `json:"contract_no"`
	SocialCategory                    string  `json:"social_category"`
	Peoples                           string  `json:"peoples"`
	CardNumber                        string  `json:"card_number"`
	HouseNumber                       string  `json:"house_number"`
	OldAddress                        string  `json:"old_address"`
	PhoneNumber1                      string  `json:"phone_number1"`
	PhoneNumber2                      string  `json:"phone_number2"`
	DateOfSigning                     string  `json:"date_of_signing"`
	DateOfDelivery                    string  `json:"date_of_delivery"`
	Proprietor                        string  `json:"proprietor"`
	Signatory                         string  `json:"signatory"`
	ChangeMethod                      string  `json:"change_method"`
	CollectiveLandPropertyCertificate string  `json:"collective_land_property_certificate"`
	Registration                      string  `json:"registration"`
	InitialHQArea                     float64 `json:"initial_hq_area"`
	RemainingHQArea                   float64 `json:"remaining_hq_area"`
	Desc                              string  `json:"desc"`
	IsDelivery                        bool    `json:"is_delivery"`
	HouseWriteOff                     bool    `json:"house_write_off"`
	CanDeclare                        bool    `json:"can_declare"`
	Single                            bool    `json:"single"`
	StagingID                         uint    `json:"staging_id"`
	DeclarationActive                 uint    `json:"declaration_active"`
	ResultCount                       uint    `json:"result_count"`
}

func UpdateHouseWriteOffStatusList(ContractNoList []string, HouseWriteOffS bool) (err error) {
	sql := db.Model(&Contract{}).Where("contract_no IN (?)", ContractNoList).Updates(map[string]interface{}{"house_write_off": HouseWriteOffS})
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

func FindContractById(id string) (*Contract, error) {
	var contract Contract
	if err := db.Where("contract_no = ?", id).First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func FindInContractById(id []string) (*Contract, error) {
	var contract Contract
	if err := db.Where("contract_no in (?)", id).Find(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func UpdateContract(contract *Contract) (err error) {
	sql := db.Model(&Contract{}).Where("contract_no = ?", contract.ContractNo).Save(&contract)
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error

}

func UpdateCardNumber(contractNo, cardNumber string) error {
	sql := db.Model(&Contract{}).Where("contract_no = ?", contractNo).Updates(map[string]interface{}{"card_number": cardNumber, "single": false})
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

type DeclarationStatusCount struct {
	DeclarationActive uint `json:"declaration_active"`
	ResultCount       uint `json:"result_count"`
}

func UpdateTargetPlacementAreaAndTemporaryRelocationArea(contractNo string, targetPlacementArea, temporaryRelocationArea float64) error {
	if contractNo == "" {
		return nil
	}
	return db.Model(&Contract{}).Where("contract_no = ?", contractNo).
		Updates(map[string]interface{}{
			"target_placement_area":     targetPlacementArea,
			"temporary_relocation_area": temporaryRelocationArea,
		}).Error
}

func FindInContract(contractNoList []string) ([]Contract, error) {
	var contracts []Contract
	err := db.Model(&Contract{}).Where("contract_no in (?)", contractNoList).Find(&contracts).Error
	return contracts, err
}

func GetLikeContracts(staginId uint, filterName string, page uint, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]Contract{},
	}
	arg1 := "%" + filterName + "%"
	return q.SearchAll(db.Model(&Contract{}).Where("staging_id = ? AND (peoples LIKE ? OR card_number LIKE ? OR Phone_number1 LIKE ? OR Phone_number2 LIKE ?) ", staginId, arg1, arg1, arg1, arg1).Order("created_at desc"))
}

type DeclarationCountList struct {
	DeclarationActive string `json:"declaration_active" form:"declaration_active"`
	ResultCount       string `json:"result_count" form:"result_count"`
}

func GetDecLareationCount(contractNo string) (*DeclarationCountList, error) {
	var count int
	err := db.Model(&Declaration{}).Where("contract_no = ?", contractNo).Count(&count).Error
	if err != nil {
		return nil, err
	}
	var count2 int
	err = db.Model(&Result{}).Joins("INNER JOIN declaration d on declaration_id = d.id").Where("d.contract_no = ?", contractNo).Count(&count2).Error
	if err != nil {
		return nil, err
	}
	test := DeclarationCountList{
		DeclarationActive: strconv.Itoa(int(count)),
		ResultCount:       strconv.Itoa(int(count2)),
	}
	return &test, nil
}

func (c *Contract) FindContractByStagingId(stagingid uint) (*Contract, error) {
	err := db.Where("staging_id = ?", stagingid).First(&c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

func UpdateCanDeclarationStatusByStaginId(stagingid uint, status int) error {
	sql := db.Model(&Declaration{}).Where("staging_id = ?", stagingid).Update("declaration_status", status)
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

func (c *Contract) UpdateContract(contractNo string) error {
	sql := db.Model(c).Where("contract_no = ?", contractNo).Update(&c)
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

func UpdateCanDeclare(contractNo []string, canDeclare bool) error {
	sql := db.Model(&Contract{}).Where("contract_no IN (?)", contractNo).Update("can_declare", canDeclare).Error
	return sql
}
