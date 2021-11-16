package v1

//TODO已完成: 根据分期期数（必须）、身份证号、合同号、手机号 模糊查询分页查询合同单列表（筛选功能）(管理员)
//TODO已完成: 根据合同号(批量)设置可否申报状态(管理员)
//TODO已完成: 根据合同号修改数据(管理员)
//TODO已完成: 根据分期期数新增合同单(管理员)
//TODO已完成: 根据合同号确定每一项合同是否完成房屋注销的选项(管理员)
import (
	"encoding/json"
	"fmt"
	"relocate/api"
	"relocate/middleware"
	"relocate/model"
	"relocate/service/contract_service"
	"relocate/util/app"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//根据分期期数（必须）、姓名、身份证号、合同号、手机号 模糊查询分页查询合同单列表（筛选功能）
func GetContractList(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.ContractFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	contractList, err := model.GetLikeContracts(body.StagingId, body.FilterName, body.Page, body.PageSize)

	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(contractList)
}

type UpdateContractNoBody struct {
	ContractNoList []string `json:"contract_no_list" form:"contract_no_list"` //要批量修改的合同号数组
	CanDeclare     bool     `json:"can_declare" form:"can_declare"`           //要修改的申报状态(是否可申报)
}

//@Tags 合同单
//@Summary 根据分期期数新增合同单
//@Description 根据分期期数新增合同单
//@Produce  json
//@Security ApiKeyAuth
//@Param data body api.AddContractBody true "合同信息" "多个身份证用英文逗号隔开"
//@Success 200 {object} app.Response
//@Failure 500 {object} app.Response
//@Router /api/v1/contract/new [post]
func NewContract(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.AddContractBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	_, err := model.FindContractById(body.ContractNo)
	if err != gorm.ErrRecordNotFound {
		return
	}
	_, err = model.FindStagingById(body.StagingID)
	if appG.HasError(err) {
		return
	}
	claim := middleware.GetClaims(c)
	data, err := json.Marshal(body)
	if appG.HasError(err) {
		return
	}
	if appG.HasError(contract_service.AddContract(body, claim.Issuer, string(data))) {
		return
	}
	appG.SuccessResponse("新增成功")
}

type UpdateHouseWriteOffBody struct {
	ContractNoList []string `json:"contract_no_list" form:"contract_no_list"` //合同单列表
	HouseWriteOff  bool     `json:"house_write_off" form:"house_write_off"`   //是否注销
}

//@Tags 合同单
//@Summary 根据合同号确定每一项合同是否完成房屋注销的选项
//@Description 传入一个合同号数组和一个要设置的房屋注销选项，批量修改
//@Produce  json
//@Security ApiKeyAuth
//@Param data body UpdateHouseWriteOffBody true "根据合同号批量修改是否注销"
//@Success 200 {object} app.Response
//@Failure 500 {object} app.Response
//@Router /api/v1/contract/updateHouseWriteOff [post]
func UpdateHouseWriteOffList(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UpdateHouseWriteOffBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	claim := middleware.GetClaims(c)
	data, err := json.Marshal(body)
	if appG.HasError(err) {
		return
	}
	if appG.HasError(contract_service.UpdateHouseWriteOffStatusList(body.ContractNoList, body.HouseWriteOff, claim.Issuer, string(data))) {
		return
	}
	appG.SuccessResponse("设置成功")
}

type ContractNoBody struct {
	ContractNo string `json:"contract_no" form:"contract_no" validate:"required"`
}

type AddCardNumberBody struct {
	CardNumber string `json:"card_number" form:"card_number" validate:"required"`
	ContractNo string `json:"contract_no" form:"contract_no" validate:"required"`
}

//@Tags 合同单
//@Summary 根据合同号增加身份证
//@Description 根据合同号增加身份证
//@Produce  json
//@Security ApiKeyAuth
//@Param data body AddCardNumberBody true "证件号"
//@Success 200 {object} app.Response
//@Failure 500 {object} app.Response
//@Router /api/v1/contract/addCardNumber [post]
func AddCardNumber(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body AddCardNumberBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	contract, err := model.FindContractById(body.ContractNo)
	if err != nil {
		appG.BadResponse("合同号不存在")
		return
	}
	if strings.Contains(contract.CardNumber, body.CardNumber) {
		appG.SuccessResponse("该身份证已添加")
		return
	}
	cardNumber := strings.Replace(body.CardNumber, "（", "(", -1)
	cardNumber = strings.Replace(cardNumber, "）", ")", -1)
	claim := middleware.GetClaims(c)
	data, err := json.Marshal(body)
	if appG.HasError(err) {
		return
	}
	if appG.HasError(contract_service.UpdateCardNumber(body.ContractNo, contract.CardNumber+","+cardNumber, claim.Issuer, string(data))) {
		return
	}
	appG.SuccessResponse("添加成功")
}

// @Tags 合同单
// @Summary 补充导入指标安置面积、计算临迁费面积
// @Accept 	multipart/form-data
// @Produce  json
// @Security ApiKeyAuth
// @Param excel formData file true "拆迁人Excel原始数据"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/contract/supplement [post]
func SupplementContract(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	fileHeader, err := api.Upload(c, "excel", 30, []string{".xlsx", ".xls"})
	if appG.HasError(err) {
		return
	}
	file, err := fileHeader.Open()
	defer file.Close()
	if appG.HasError(err) {
		return
	}
	count, err := contract_service.SupplementContract(file)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(fmt.Sprintf("导入成功,条数为:%d", count))
}

//
func GetDeclarationCount(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body ContractNoBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	contract, err := model.GetDecLareationCount(body.ContractNo)
	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse(contract)
}

type setCanDeclareBody struct {
	StagingId uint `json:"staging_id" form:"staging_id" validate:"required"`
	Status    int  `json:"status" form:"status" validate:"oneof=0 1"`
}

//根据分期id，设置该分期的所有合同号的申报状态
func SettingCanDeclareByStaginId(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body setCanDeclareBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	if appG.HasError(contract_service.SettingCanDeclareByStagingId(body.StagingId, body.Status)) {
		return
	}
	appG.SuccessResponse("设置申报状态成功")

}

type UpdataContractBody struct {
	CanDeclare                        bool   `json:"can_declare" form:"can_declare"`
	ChangeMethod                      string `json:"change_method" form:"change_method"`
	CollectiveLandPropertyCertificate string `json:"collective_land_property_certificate" form:"collective_land_property_certificate"`
	ContractNo                        string `json:"contract_no" form:"contract_no" validate:"required"`
	DateOfDelivery                    string `json:"date_of_delivery" form:"date_of_delivery" validate:"required"`
	DateOfSigning                     string `json:"date_of_signing" form:"date_of_signing" validate:"required"`
	Desc                              string `json:"desc" form:"desc"`
	HouseNumber                       string `json:"house_number" form:"house_number" validate:"required"`
	HouseWriteOff                     bool   `json:"house_write_off" form:"house_write_off"`
	OldAddress                        string `json:"old_address" form:"old_address" validate:"required"`
	PhoneNumber1                      string `json:"phone_number1" form:"phone_number1"`
	PhoneNumber2                      string `json:"phone_number2" form:"phone_number2"`
	Proprietor                        string `json:"proprietor" form:"proprietor"`
	Registration                      string `json:"registration" form:"registration"`
	Signatory                         string `json:"signatory" form:"signatory" validate:"required"`
	SocialCategory                    string `json:"social_category" form:"social_category"`
	StagingID                         uint   `json:"staging_id" form:"staging_id" validate:"required"`
}

//根据合同号修改数据
func UpdateContract(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UpdataContractBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(contract_service.UpdateContract(body.CanDeclare, body.ChangeMethod,
		body.CollectiveLandPropertyCertificate, body.ContractNo, body.DateOfDelivery,
		body.DateOfSigning, body.Desc, body.HouseNumber, body.HouseWriteOff, body.OldAddress, body.PhoneNumber1,
		body.PhoneNumber2, body.Proprietor, body.Registration, body.Signatory, body.SocialCategory, body.StagingID)) {
		return
	}

	appG.SuccessResponse("修改合同单成功")
}

type UpdateCanDeclareBody struct {
	CanDeclare bool     `json:"can_declare" from:"can_declare"`
	Contractno []string `json:"contract_no_list" form:"contract_no_list"`
}

func UpdateCanDeclare(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UpdateCanDeclareBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(contract_service.UpdateCanDeclare(body.Contractno, body.CanDeclare)) {
		return
	}
	appG.SuccessResponse("（批量）修改状态成功")
}
