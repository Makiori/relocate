package router

import (
	v1 "relocate/api/v1"
	_ "relocate/docs"
	"relocate/middleware"
	"relocate/util/sign"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

//说明
//middleware.JWT()：管理员、用户共有
//middleware.JWT(sign.AdminClaimsType)：管理员独有
//middleware.JWT(sign.UserClaimsType)：用户独有

//初始化路由信息
func InitRouter() *gin.Engine {
	r := gin.New()
	//全局 Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	r.Use(gin.Recovery())
	//全局 日志中间件
	r.Use(middleware.LoggerToFile())
	//全局 跨域中间件
	r.Use(middleware.Cors())
	//加载模板文件
	r.LoadHTMLGlob("router/templates/*")
	//加载静态文件
	r.Static("/web", "router/static")
	//swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//v1版本
	apiV1 := r.Group("/api/v1")
	initAdminRouter(apiV1)
	initConfigRouter(apiV1)
	initStagingRouter(apiV1)
	initContractRouter(apiV1)
	initUserRouter(apiV1)
	initTimeRouter(apiV1)
	initHuxingRouter(apiV1)
	initInformationRouter(apiV1)
	initDeclarationRouter(apiV1)
	initAreaDetailRouter(apiV1)
	initResultRouter(apiV1)
	initLoggingRouter(apiV1)
	initAccountingRouter(apiV1)
	initCheckRouter(apiV1)
	return r
}

func initAdminRouter(apiV1 *gin.RouterGroup) {
	admin := apiV1.Group("/admin")
	{
		admin.POST("/login", v1.AdminLogin)
		admin.POST("/update", middleware.JWT(sign.AdminClaimsType), v1.AdminUpdatePassword)
		admin.GET("/get", middleware.JWT(sign.AdminClaimsType), v1.GetAllAdmin)
	}
}

func initConfigRouter(apiV1 *gin.RouterGroup) {
	config := apiV1.Group("/config")
	{
		config.GET("/staging/setting", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingStagingConfig)
		config.GET("/rounds/setting", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingNowRoundsConfig)
		config.GET("/staging/get", v1.GetStagingConfig)
		config.GET("/rounds/get", v1.GetNowRoundsConfig)
		config.POST("/huxing/groupingOptional/setting", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingHuxingGroupingOptionalConfig)
		config.GET("/huxing/groupingOptional/get", v1.GetHuxingGroupingOptionalConfig)
		config.GET("/declaration/active/setting", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingIncludeTotalConfig)
		config.GET("/declaration/active/get", v1.GetIncludeTotalConfig)
		config.GET("/huxing/optional/get", v1.GetHuxingOptionalConfig)
	}
}

func initStagingRouter(apiV1 *gin.RouterGroup) {
	staging := apiV1.Group("/staging")
	{
		staging.POST("/new", middleware.JWT(sign.AdminClaimsType), v1.NewStaging)
		staging.GET("/get", v1.GetStaging)
		staging.POST("/import", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.ImportContractByStaging)
		staging.GET("/getCount", middleware.JWT(sign.AdminClaimsType), v1.GetStagingContractCount)
	}
}

func initContractRouter(apiV1 *gin.RouterGroup) {
	contract := apiV1.Group("/contract")
	{
		contract.GET("/get", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.GetContractList)
		contract.GET("/getDeclarationCount", middleware.JWT(), v1.GetDeclarationCount)
		contract.POST("/updateHouseWriteOff", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateHouseWriteOffList)
		contract.POST("/new", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.NewContract)
		contract.GET("/setCanDeclareByStagingId", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingCanDeclareByStaginId)
		contract.POST("/addCardNumber", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.AddCardNumber)
		contract.POST("/supplement", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SupplementContract)
		contract.POST("/update", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateContract)
		contract.POST("/updateCanDeclare", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateCanDeclare)

	}
}

func initUserRouter(apiV1 *gin.RouterGroup) {
	user := apiV1.Group("/user")
	{
		user.POST("/login", v1.UserLogin)
		user.GET("/findMyself", middleware.JWT(sign.UserClaimsType), v1.UserQuery)
		user.GET("/get", v1.UserRegisterSendCode)
		user.GET("/forget", v1.UserForgetSendCode)
		user.POST("/register", v1.UserRegister)
		user.POST("/validator", middleware.JWT(sign.UserClaimsType), v1.UserValidator)
		user.GET("/getAllValidate", middleware.JWT(sign.UserClaimsType), v1.GetAllValidate)
		user.POST("/updateStatus", middleware.JWT(sign.AdminClaimsType), v1.UpdateStatus)
		user.GET("/getContract", middleware.JWT(sign.UserClaimsType), v1.GetContractNoByCardNumber)
		user.POST("/getResult", middleware.JWT(sign.UserClaimsType), v1.GetResultByCardNumber)

	}
}
func initTimeRouter(apiV1 *gin.RouterGroup) {
	/*time := apiV1.Group("/time")
	{
	}*/
}

func initHuxingRouter(apiV1 *gin.RouterGroup) {
	huxing := apiV1.Group("/huxing")
	{
		huxing.GET("/get", v1.GetHuxing)
		huxing.GET("/get/optional", v1.GetOptionalHuxing)
		huxing.POST("/new", middleware.JWT(sign.AdminClaimsType), v1.NewHuxing)
		huxing.POST("/delete", middleware.JWT(sign.AdminClaimsType), v1.DeleteHuxing)
		huxing.POST("/update", middleware.JWT(sign.AdminClaimsType), v1.UpdateHuxing)
	}
}

func initInformationRouter(apiV1 *gin.RouterGroup) {
	information := apiV1.Group("/information")
	{
		information.POST("/upload", v1.UploadImage)
		information.GET("/get", middleware.JWT(sign.AdminClaimsType), v1.GetInformationList)
		information.GET("/getpublic", v1.GetPublicInformationList)
		information.POST("/new", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SetInformation)
		information.POST("/updateStatus", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateInformationStatus)
		information.POST("/delete", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.DeleteInformation)
		information.POST("/update", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateInformation)
	}
}

func initDeclarationRouter(apiV1 *gin.RouterGroup) {
	declaration := apiV1.Group("/declaration")
	{
		declaration.POST("/add", middleware.JWT(), v1.AddDeclaration)
		declaration.POST("/addNew", middleware.JWT(), v1.AddDeclarationNew)
		declaration.POST("/delete", middleware.JWT(sign.AdminClaimsType), v1.DeleteDeclarationResult)
		declaration.POST("/enterResult", middleware.JWT(sign.AdminClaimsType), v1.EnterResult)
		declaration.GET("/getAdminName", middleware.JWT(sign.AdminClaimsType), v1.GetAdminName)
		declaration.POST("/updateDeclarationStatus", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateDeclarationStatus)
		declaration.POST("/printing", middleware.JWT(sign.AdminClaimsType), v1.Printing)
		declaration.GET("/get", middleware.JWT(sign.AdminClaimsType), v1.GetDeclarationList)
		declaration.GET("/getAll", middleware.JWT(sign.AdminClaimsType), v1.GetAllDeclaration)
		declaration.GET("/getDeclaration", middleware.JWT(sign.AdminClaimsType), v1.GetDeclarationByContractNo)
		declaration.POST("/updateDeclaration", v1.UpdateDeclaration)
		declaration.POST("/updateActive", middleware.JWT(sign.AdminClaimsType), v1.UpdateDeclarationActive)
		declaration.GET("/getDetail", v1.GetDeclarationDetail)
		declaration.GET("/export", v1.ExportDeclaration)
	}
}

func initAreaDetailRouter(apiV1 *gin.RouterGroup) {
	areaDetail := apiV1.Group("/areaDetail")
	{
		areaDetail.GET("/get", middleware.JWT(), v1.GetAreaDetailList)
	}
}

func initResultRouter(apiV1 *gin.RouterGroup) {
	result := apiV1.Group("/result")
	{
		result.GET("/get", v1.GetResultList)
		result.POST("/export", middleware.JWT(sign.AdminClaimsType), v1.ExportResults)
		result.GET("/getByStatus", v1.GetResultByStatusList)
		result.POST("/updatePublicityStat", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdatePublicityStat)
	}
}

func initLoggingRouter(apiV1 *gin.RouterGroup) {
	logging := apiV1.Group("/logging")
	{
		logging.GET("/get", middleware.JWT(sign.AdminClaimsType), v1.GetLogging)
	}
}
func initAccountingRouter(apiV1 *gin.RouterGroup) {
	accounting := apiV1.Group("/accounting")
	{
		accounting.POST("", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.AddAccounting)
		accounting.GET("", middleware.JWT(sign.AdminClaimsType), v1.GetAccountingList)
		accounting.GET("/export", middleware.JWT(sign.AdminClaimsType), v1.ExportAccounting)

	}
}

func initCheckRouter(apiV1 *gin.RouterGroup) {
	check := apiV1.Group("/check")
	{
		check.GET("", middleware.JWT(sign.AdminClaimsType), v1.GetCheckList)
		check.POST("", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.AddCheck)
		check.GET("/export", middleware.JWT(sign.AdminClaimsType), v1.ExportCheck)

	}
}
