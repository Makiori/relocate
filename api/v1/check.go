package v1

import (
	"fmt"
	"net/url"
	"relocate/api"
	"relocate/model"
	"relocate/service/check_service"
	"relocate/util/app"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags 核算
// @Summary 根据合同号、被拆迁人模糊查询分页查询核算列表（筛选功能）(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param filterName query string false "搜索关键字(合同号、被拆迁人)"
// @Param page query int false "页码"
// @Param pageSize query int false "页面大小"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/check [get]
func GetCheckList(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.CheckFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	data, err := model.GetLikeCheck(body)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}

type AddCheckBody struct {
	ResultID    uint    `json:"result_id" validate:"required"`
	RealityArea float64 `json:"reality_area"`
}

func AddCheck(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body AddCheckBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(check_service.AddCheck(body.ResultID, body.RealityArea)) {
		return
	}
	appG.SuccessResponse("核算成功")
}

func ExportCheck(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.CheckFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	body.Page = 1
	body.PageSize = 9999
	data, err := model.GetLikeCheck(body)
	if appG.HasError(err) {
		return
	}
	checkDataList, ok := data.Data.(*[]model.Check)
	if !ok {
		appG.BadResponse("导出过程发生异常")
		return
	}
	file, filename, err := check_service.ExportExcel(*checkDataList)
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
