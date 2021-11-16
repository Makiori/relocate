package v1

import (
	"fmt"
	"net/url"
	"relocate/api"
	"relocate/middleware"
	"relocate/model"
	"relocate/service/declaration_service"
	"relocate/util/app"
	"relocate/util/sign"
	"strconv"

	"github.com/gin-gonic/gin"
)

//TODO已完成: 根据分期期数（必须）、身份证号、合同号、手机号 模糊查询分页查询申报表列表（筛选功能）(管理员)
//TODO已完成: 根据合同号查看申报表详情
//TODO已完成: 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
//TODO已完成: 根据合同号更改申报状态(管理员)
//TODO已完成: 申报表打印管理员姓名(管理员)
//TODO已完成: 根据合同号更改中签状态（是否录入结果）(管理员)

// @Tags 申报单
// @Summary 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
// @Description 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
// @Produce  json
// @Security ApiKeyAuth
// @Param data body api.AddDeclarationBody true "申报信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/addNew [post]
func AddDeclarationNew(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.AddDeclarationBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	//判断是管理员登录还是用户登录 -- true：管理员 false：用户
	claim := middleware.GetClaims(c)
	ok := false
	if claim.Type == sign.AdminClaimsType {
		ok = true
	}
	if appG.HasError(declaration_service.AddDeclarationNew(body, claim.Issuer, ok)) {
		return
	}
	appG.SuccessResponse("申报成功")
}

// @Tags 申报单
// @Summary 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
// @Description 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
// @Produce  json
// @Security ApiKeyAuth
// @Param data body api.AddDeclarationBody true "申报信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/add [post]
func AddDeclaration(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.AddDeclarationBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	//判断是管理员登录还是用户登录 -- true：管理员 false：用户
	claim := middleware.GetClaims(c)
	ok := false
	if claim.Type == sign.AdminClaimsType {
		ok = true
	}
	if appG.HasError(declaration_service.AddDeclarationNew(body, claim.Issuer, ok)) {
		return
	}
	appG.SuccessResponse("申报成功")
}

// @Tags 申报单
// @Summary 获取所有的申报表列表
// @Description 获取所有的申报表列表
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/getAll [get]
func GetAllDeclaration(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	data, err := model.GetAllDeclaration()
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}

func GetDeclarationList(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.DeclarationFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	declaration, err := model.GetDeclaration(body.StagingId, body.HuxingId, body.TimeId, body.DeclarationStatus, body.WinningStatus, body.ActiveState, body.FilterName, body.Page, body.PageSize)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(declaration)
}

type QueryDeclarationBody struct {
	ContractNo    string `json:"contract_no" form:"contract_no" validate:"required"`
	StagingId     string `json:"staging_id" form:"staging_id" validate:"required"`
	DeclarationID int    `json:"declaration_id" form:"declaration_id" validate:"required"`
}

func GetDeclarationDetail(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body QueryDeclarationBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	declaration, err := model.GetDeclarationDetail(body.ContractNo, body.StagingId, body.DeclarationID)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(declaration)
}

// @Tags 申报单
// @Summary 申报表打印管理员姓名(管理员)
// @Description 申报表打印管理员姓名(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/getAdminName [get]
func GetAdminName(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	claim := middleware.GetClaims(c)
	admin, err := model.GetAdminInfo(claim.Issuer)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(map[string]interface{}{
		"name": admin.AdminName,
	})
}

type UpdateStatusBody struct {
	DeclarationID uint `json:"declaration_id" form:"declaration_id"`
	Status        int  `json:"status" form:"status" validate:"oneof=0 1"`
}

// @Tags 申报单
// @Summary 根据申报表ID更改申报状态(管理员)
// @Description 根据申报表ID更改申报状态(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param declaration_id query int true "申报ID（必须）"
// @Param status query int true "申报表申报状态（必须）-- 0：表示进行中；1：表示已确定"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/updateDeclarationStatus [post]
func UpdateDeclarationStatus(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UpdateStatusBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	claim := middleware.GetClaims(c)
	if appG.HasError(declaration_service.ChangeDeclarationStatus(body.DeclarationID, body.Status, claim.Issuer)) {
		return
	}
	appG.SuccessResponse("修改成功")
}

type UpdateDeclarationBody struct {
	DeclarationID uint
}

// @Tags 申报单
// @Summary 根据申报表ID更改申报表数据
// @Description 根据申报表ID更改申报表数据
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/updateDeclaration [post]
func UpdateDeclaration(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.UpdateDeclaration
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(declaration_service.UpdateDeclaration(body.DeclarationID, body.HuxingID, body.Trustee, body.TrusteePhoneNumber, body.TrusteeCardNumber, body.TrusteeRelationship)) {
		return
	}
	appG.SuccessResponse("修改申报表数据成功")
}

type PrintingBody struct {
	DeclarationID uint `json:"declaration_id" form:"declaration_id" validate:"required"`
}

// @Tags 申报单
// @Summary 打印申报表(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param declaration_id query int true "申报表ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/printing [post]
func Printing(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body PrintingBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	claim := middleware.GetClaims(c)
	admin, err := model.GetAdminInfo(claim.Issuer)
	if appG.HasError(err) {
		return
	}
	file, filename, err := declaration_service.GenerateExcelNew(body.DeclarationID, admin.AdminSignname)
	if appG.HasError(err) {
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(filename)))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	buffer, err := file.WriteToBuffer()
	if appG.HasError(err) {
		return
	}
	c.Writer.Header().Add("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	if appG.HasError(file.Write(c.Writer)) {
		return
	}
}

type DeleteResultBody struct {
	DeclarationID uint `json:"declaration_id" form:"declaration_id" validate:"required"`
}

type UpdateDeclarationActiveBody struct {
	DeclarationID uint `json:"declaration_id" form:"declaration_id" validate:"required"`
	State         bool `json:"state" form:"state" validate:"omitempty,required"`
}

// @Tags 申报单
// @Summary 修改申报有效状态，根据申报ID
// @Description 修改申报有效状态，根据申报ID
// @Produce  json
// @Security ApiKeyAuth
// @Param data body api.AddDeclarationBody true "申报信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/updateActive [post]

func UpdateDeclarationActive(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UpdateDeclarationActiveBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	if appG.HasError(declaration_service.UpdateDeclarationActive(body.DeclarationID, body.State)) {
		return
	}
	appG.SuccessResponse("修改申报有效状态成功")
}

// @Tags 申报单
// @Summary 删除申报结果，根据申报ID清除
// @Description 删除申报结果，根据申报ID清除
// @Produce  json
// @Security ApiKeyAuth
// @Param data body api.AddDeclarationBody true "申报信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/delete [post]
func DeleteDeclarationResult(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body DeleteResultBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	if appG.HasError(declaration_service.DeleteDeclarationResult(body.DeclarationID)) {
		return
	}
	appG.SuccessResponse("删除申报结果成功")
}

type EnterResultBody struct {
	BuildingNo    string `json:"building_no" from:"building_no"`
	DeclarationID uint   `json:"declaration_id" form:"declaration_id" validate:"required"`
	RoomNo        string `json:"room_no" form:"room_no"`
	Status        int    `json:"status" form:"status" validate:"oneof=0 1"`
}

func EnterResult(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body EnterResultBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	if appG.HasError(declaration_service.EnterResult(body.DeclarationID, body.BuildingNo, body.RoomNo, body.Status)) {
		return
	}
	appG.SuccessResponse("更改中签状态成功")
}

type QueryDeclarationByContractNoBody struct {
	ContractNo string `json:"contract_no" form:"contract_no" validate:"required"`
}

func GetDeclarationByContractNo(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body QueryDeclarationByContractNoBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	DeclarationList, err := model.FindAllDeclaration(body.ContractNo)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(DeclarationList)
}

/*func ExportDeclaration(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.DeclarationFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	body.PageSize = 9999
	data, err := model.GetLikeResults(body)
	if appG.HasError(err) {
		return
	}
	resultDataList, ok := data.Data.(*[]model.ResultData)
	if !ok {
		appG.BadResponse("导出过程发生异常")
		return
	}
	file, filename, err := result_service.ExportExcel(*resultDataList)
	if appG.HasError(err) {
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(filename)))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	buffer, err := file.WriteToBuffer()
	if appG.HasError(err) {
		return
	}
	c.Writer.Header().Add("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	if appG.HasError(file.Write(c.Writer)) {
		return
	}
}*/
