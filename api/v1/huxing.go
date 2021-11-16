package v1

//TODO已完成: 获取户型数据列表
//TODO已完成: 根据ID删除户型数据(管理员)
//TODO已完成: 根据ID修改户型数据(管理员)
//TODO已完成: 新增户型数据(管理员)
import (
	"relocate/model"
	"relocate/service/huxing_service"
	"relocate/util/app"

	"github.com/gin-gonic/gin"
)

// @Tags 户型
// @Summary 获取户型
// @Description 获取户型
// @Param staging_id query int true "分期期数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/huxing/get/optional [get]

type GetHuxingListBody struct {
	StagingID uint `json:"staging_id" form:"staging_id" validate:"required"`
}

func GetHuxing(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body GetHuxingListBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	huxingList, err := model.FindAllHuxing(body.StagingID)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(huxingList)
}

// @Tags 户型
// @Summary 获取户型（可选>0以上的）
// @Description 获取户型（可选>0以上的）
// @Produce  json
// @Param staging_id query int true "分期期数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/huxing/get/optional [get]
func GetOptionalHuxing(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body GetHuxingListBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	huxingList, err := model.FindAllOptionalHuxing(body.StagingID)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(huxingList)
}

// @Tags 户型
// @Summary 新增户型
// @Description 新增户型
// @Produce  json
// @Param data body
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/huxing/new [post]
type NewHuxingBody struct {
	Area       string `json:"area" form:"area" validate:"required"`
	AreaShow   string `json:"area_show" form:"area_show" validate:"required"`
	BuildingNo string `json:"building_no" form:"building_no" validate:"required"`
	HuxingNo   string `json:"huxing_no" form:"huxing_no" validate:"required"`
	Maximum    uint   `json:"maximum" form:"maximum"`
	Quantity   uint   `json:"quantity" form:"quantity" validate:"required"`
	Rounds     uint   `json:"rounds" form:"rounds"`
	StagingId  uint   `json:"staging_id" form:"staging_id"`
}

func NewHuxing(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body NewHuxingBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(huxing_service.AddHuxing(body.Area, body.AreaShow, body.BuildingNo, body.HuxingNo, body.Maximum, body.Quantity, body.Rounds, body.StagingId)) {
		return
	}

	appG.SuccessResponse("新增户型成功")
}

// @Tags 户型
// @Summary 删除户型
// @Description 删除户型
// @Produce  json
// @Param data body
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/huxing/delete [post]
type DeleteHuxingBody struct {
	Id uint `json:"id" form:"id" validate:"required"`
}

func DeleteHuxing(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body DeleteHuxingBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(huxing_service.DeleteHuxing(body.Id)) {
		return
	}
	appG.SuccessResponse("删除户型成功")
}

// @Tags 户型
// @Summary 修改户型
// @Description 修改户型
// @Produce  json
// @Param data body
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/huxing/update [post]
type UpdateHuxingBody struct {
	Area       string `json:"area" form:"area" validate:"required"`
	AreaShow   string `json:"area_show" form:"area_show" validate:"required"`
	BuildingNo string `json:"building_no" form:"building_no"`
	HuxingNo   string `json:"huxing_no" form:"huxing_no"`
	Maximum    uint   `json:"maximum" form:"maximum"`
	Quantity   uint   `json:"quantity" form:"quantity"`
	Rounds     uint   `json:"rounds" form:"rounds"`
	Id         uint   `json:"id" form:"id" validate:"required"`
}

func UpdateHuxing(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body UpdateHuxingBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(huxing_service.UpdateHuxing(body.Area, body.AreaShow, body.BuildingNo, body.HuxingNo, body.Maximum, body.Quantity, body.Rounds, body.Id)) {
		return
	}
	appG.SuccessResponse("修改户型成功")
}
