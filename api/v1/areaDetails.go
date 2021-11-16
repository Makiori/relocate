package v1

import (
	"relocate/api"
	"relocate/model"
	"relocate/util/app"

	"github.com/gin-gonic/gin"
)

//TODO已完成: 分页获取合同号面积明细记录
//TODO: 用户注册、申报、交易等都要进行记录明细

func GetAreaDetailList(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.AreaDetailFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	areaDetailList, err := model.GetAreaDetail(body.ContractNo, body.Page, body.PageSize)

	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(areaDetailList)
}
