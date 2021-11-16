package v1

import (
	"relocate/model"
	"relocate/util/app"

	"github.com/gin-gonic/gin"
)

//TODO已完成: 根据操作人(必须)、时间区间(非必须)分页获取日志数据列表(管理员)
//TODO已完成: 采用中间件记录操作日志

/*// @Tags 日志
// @Summary 分页获取所有的日志数据列表
// @Description 分页获取所有的日志数据列表
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/logging/getAll [get]
func GetAllLogging(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	data, err := model.GetAllLogging()
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}
*/

// @Tags 日志
// @Summary 分页获取所有的日志数据列表
// @Description 分页获取所有的日志数据列表
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/logging/getAll [get]
func GetLogging(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	data, err := model.GetLogging()
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}
