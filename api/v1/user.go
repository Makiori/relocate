package v1

import (
	"relocate/model"
	"relocate/service/user_service"
	"relocate/util/app"
	"relocate/util/gredis"
	"relocate/util/random"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//TODO已完成: 用户通过手机验证码注册
//TODO已完成: 用户通过手机号码+密码登录
//TODO已完成: 获取用户名下合同号(用户合同表)
//TODO已完成: 用户进行身份核验
//TODO已完成: 用户通过验证码修改密码

type UserGetCodeBody struct {
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,checkMobile,len=11"`
}

// @Tags 用户
// @Summary 获取验证码(注册)
// @Produce json
// @Param phone_number query string true "手机号码"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/get [get]
func UserRegisterSendCode(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserGetCodeBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	_, err := model.GetUserInfo(body.PhoneNumber)
	if err == nil {
		appG.BadResponse("用户已存在")
		return
	}
	code := random.Code(6)
	if !random.SendCode(body.PhoneNumber, code) {
		appG.BadResponse("发送验证码错误")
		return
	}
	user := make(map[string]string)
	user["user-register-"+body.PhoneNumber] = code
	gredis.Set(user, 200)
	appG.SuccessResponse("发送验证码成功")
}

type UserRegisterBody struct {
	Username string `json:"username" validate:"required,len=11,checkMobile" `
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required,len=6"` //验证码为6位数字
}

func UserForgetSendCode(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserGetCodeBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	_, err := model.GetUserInfo(body.PhoneNumber)
	if err != nil {
		appG.BadResponse("用户不存在")
		return
	}
	code := random.Code(6)
	if !random.SendCode(body.PhoneNumber, code) {
		appG.BadResponse("发送验证码错误")
		return
	}
	user := make(map[string]string)
	user["user-register-"+body.PhoneNumber] = code
	gredis.Set(user, 200)
	appG.SuccessResponse("发送验证码成功")

}

// @Tags 用户
// @Summary 注册(通过手机号码+密码注册+验证码)
// @Produce json
// @Param data body UserRegisterBody true "注册信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/register [post]
func UserRegister(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserRegisterBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	_, err := model.GetUserInfo(body.Username)
	if err == nil {
		appG.BadResponse("用户已存在")
		return
	}
	code, _ := gredis.Get("user-register-" + body.Username)
	if code != body.Code { //用户没有获取验证码或者验证码过期
		appG.BadResponse("验证码错误")
		return
	}
	if appG.HasError(user_service.CreateUser(body.Username, body.Password)) {
		return
	}
	gredis.Delete("user-register-" + body.Username)
	appG.SuccessResponse("注册成功")
}

type UserContractBody struct {
	CardNumber string `json:"card_number" form:"card_number"`
}

type UserInfoBody struct {
	PhoneNumber   string `json:"phone_number" form:"phone_number" validate:"required,checkMobile,len=11"` //当前登录用户的手机号码
	Peoples       string `json:"peoples" form:"peoples" validate:"required"`                              //真实姓名
	IdNumberType  int    `json:"id_number_type" form:"id_number_type"`                                    //证件类型 0大陆身份证 1香港身份证 2护照
	CardNumber    string `json:"card_number" form:"card_number" validate:"required"`                      //身份证
	PositiveImage string `json:"positive_image" form:"positive_image" validate:"required"`                //身份证正面,要求上传的图片的base64格式
	NegativeImage string `json:"negative_image" form:"negative_image" validate:"required"`                //身份证反面,要求上传的图片的base64格式
	SuffixA       string `json:"suffix_a" form:"suffix_a" validate:"required"`                            //正面图片原先的后缀(例：.jpg .png等)
	SuffixB       string `json:"suffix_b" form:"suffix_b" validate:"required"`                            //背面图片原先的后缀(例：.jpg .png等)
}

// @Tags 用户
// @Summary 用户进行身份核验
// @Produce json
// @Security ApiKeyAuth
// @Param data body UserInfoBody true "用户身份信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/validator [post]
func UserValidator(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserInfoBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	cardNumber := strings.Replace(body.CardNumber, "（", "(", -1)
	cardNumber = strings.Replace(cardNumber, "）", ")", -1)
	//更新
	if appG.HasError(user_service.UpdateUser(body.IdNumberType, body.PhoneNumber, body.Peoples, cardNumber, body.PositiveImage, body.NegativeImage, body.SuffixA, body.SuffixB)) {
		return
	}
	//自动匹配
	if model.AutoMatchUser(body.PhoneNumber, body.Peoples, cardNumber) {
		//修改状态
		if appG.HasError(user_service.UpdateStatus(body.PhoneNumber)) {
			return
		}
	}
	appG.SuccessResponse("提交成功")
}

type UserBody struct {
	Phone string `json:"phone" form:"phone" validate:"required,checkMobile,len=11"`
}

// @Tags 用户
// @Summary 用户通过手机查看自己的身份信息
// @Produce json
// @Security ApiKeyAuth
// @Param phone query string true "用户通过手机号查看身份信息"  前提是通过核验
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/findMyself [get]
func UserQuery(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	date, err := model.GetUserInfo(body.Phone)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.BadResponse("账号不存在")
			return
		}
		appG.HasError(err)
		return
	}
	appG.SuccessResponse(map[string]interface{}{"user": date})
}

type UpdateStatusByPhoneNumber struct {
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,checkMobile,len=11"`
	Pass        int    `json:"pass" form:"pass" validate:"oneof=0 1 2 3"`
}

func UpdateStatus(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UpdateStatusByPhoneNumber
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(user_service.UpdateStatusByPhoneNumber(body.PhoneNumber, body.Pass)) {
		return
	}
	appG.SuccessResponse("修改用户校验状态成功")
}

type UserLoginBody struct {
	Username string `json:"username" validate:"required,checkMobile,len=11"`
	Password string `json:"password" validate:"required"`
}

// @Tags 用户
// @Summary 登录(通过手机号码+密码登录)
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/login [post]
func UserLogin(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserLoginBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	token, err := user_service.GenerateToken(body.Username, body.Password)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(token)
}

type GetAllValidateBody struct {
	FilterStatus string `json:"filterStatus" form:"filterStatus"`
	FilterName   string `json:"filterName" form:"filterName"`
	Page         uint   `json:"page" form:"page"`
	PageSize     uint   `json:"pageSize" form:"pageSize"`
}

// @Tags 用户
// @Summary 根据用户状态查看核验列表
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/getAllValidate [get]
func GetAllValidate(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body GetAllValidateBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	userList, err := model.FindLikeUser(body.FilterStatus, body.FilterName, body.Page, body.PageSize)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(userList)
}

type GetContractByCardNumberBody struct {
	CardNumber string `json:"card_number" form:"card_number"`
}

type GetResultByCardNumberBody struct {
	CardNumber string `json:"card_number" form:"card_number"`
}

func GetContractNoByCardNumber(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body GetContractByCardNumberBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	contractList, err := model.FindContractNoByCardnumber(body.CardNumber)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(contractList)
}

func GetResultByCardNumber(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body GetResultByCardNumberBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	List, err := model.FindResultByCardnumber(body.CardNumber)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(List)
}
