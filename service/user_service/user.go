package user_service

import (
	"relocate/api"
	"relocate/model"
	"relocate/util/errors"
	"relocate/util/logging"
	"relocate/util/sign"
	"relocate/util/times"

	"github.com/jinzhu/gorm"
)

//生成token
func GenerateToken(username, password string) (interface{}, error) {
	user, err := model.GetUserInfo(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.BadError("用户账号不存在")
		}
		return "", err
	}
	if user.Password != sign.EncodeMD5(password+user.Salt) {
		return "", errors.BadError("密码错误")
	}

	token, err := sign.GenerateToken(string(user.PhoneNumber), username, sign.UserClaimsType)
	if err != nil {
		return "", err
	}
	return map[string]interface{}{"token": token, "phoneNumber": user.PhoneNumber, "name": user.Name, "idNumber": user.IdNumber, "userStatus": user.UserStatus}, nil
}

func CreateUser(username, password string) (err error) {
	user, err := model.GetUserInfo(username)
	salt := "ABCDEF"
	if user != nil {
		return errors.BadError("手机号码已被使用")
	}
	user = &model.User{
		PhoneNumber: username,
		Password:    sign.EncodeMD5(password + salt),
		Salt:        salt,
		UserStatus:  model.NotVerifiedUser,
		CreatedAt:   times.JsonTime{},
	}
	return user.Create()
}

func UpdateUser(idCardType int, phoneNumber, name, idCard, imagesA, imagesB, suffixA, suffixB string) error {
	idtype := model.MainlandId
	if idCardType == 1 {
		idtype = model.HongKongId
	} else if idCardType == 2 {
		idtype = model.HuZhaoId
	}
	formattedA, err := api.UploadImg(imagesA, suffixA, phoneNumber, "positive_images")
	if err != nil {
		logging.Info(err)
		return errors.BadError("正面照上传出错")
	}
	formattedB, err := api.UploadImg(imagesB, suffixB, phoneNumber, "negative_images")
	if err != nil {
		logging.Info(err)
		return errors.BadError("反面照上传出错")
	}
	user := &model.User{
		IdNumberType:  idtype,
		PhoneNumber:   phoneNumber,
		Name:          name,
		IdNumber:      idCard,
		PositiveImage: formattedA,
		NegativeImage: formattedB,
	}
	return user.UpdateUser()
}

func UpdateStatus(phoneNumber string) error {
	user := &model.User{
		PhoneNumber: phoneNumber,
		UserStatus:  model.AutomaticMatchingUser,
	}
	return user.UpdateUser()
}

func UpdateStatusByPhoneNumber(phoneNumber string, pass int) error {
	if _, err := getUserByPhoneNumber(phoneNumber); err != nil {
		return err
	}
	if err := model.UpdateStatusByPhoneNumber(phoneNumber, pass); err != nil {
		return errors.BadError("修改用户校验状态失败")
	}
	return nil
}

func getUserByPhoneNumber(phoneNumber string) (interface{}, error) {
	user, err := model.QueryUserByPhone(phoneNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("该用户不存在")
		}
		return nil, err
	}
	return user, nil
}
