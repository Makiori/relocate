package information_server

import (
	"relocate/api"
	"relocate/model"
	"relocate/util/errors"
	"relocate/util/logging"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

func AddInformation(title string, content string, picture string, suffix string) error {
	information := model.Information{
		Title:   title,
		Content: content,
	}
	if picture != "" {
		formatTime := strconv.FormatInt(time.Now().Unix(), 10)
		formatImg, err := api.UploadImg(picture, suffix, "public", formatTime)
		if err != nil {
			logging.Info(err)
			return errors.BadError("图片上传失败")
		}
		information.Picture = formatImg
	}
	if err := information.Create(); err != nil {
		return errors.BadError("新增失败")
	}
	return nil
}

func getInformation(id uint) (*model.Information, error) {
	information, err := model.FindInformationByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("资讯不存在")
		}
		return nil, err
	}
	return information, nil
}

func DeleteInformation(id uint) error {
	information := model.Information{
		Model: model.Model{ID: id},
	}
	//判断资讯是否存在
	if _, err := getInformation(id); err != nil {
		return err
	}
	if err := information.Delete(); err != nil {
		return errors.BadError("删除失败")
	}
	return nil
}

func UpdateInformationStatus(id uint, status model.InformationStatus) error {

	//判断资讯是否存在
	i, err := getInformation(id)
	if err != nil {
		return err
	}
	i.Status = status
	if err := i.UpdateStatus(); err != nil {
		return errors.BadError("更新状态失败")
	}
	return nil
}

func UpdateInformation(id uint, title string, content string, picture string, suffix string) error {
	information := model.Information{
		Model:   model.Model{ID: id},
		Title:   title,
		Content: content,
	}
	if picture != "" {
		formatTime := strconv.FormatInt(time.Now().Unix(), 10)
		formatImg, err := api.UploadImg(picture, suffix, "public", formatTime)
		if err != nil {
			logging.Info(err)
			return errors.BadError("图片上传失败")
		}
		information.Picture = formatImg
	}
	if _, err := getInformation(id); err != nil {
		return err
	}
	if err := information.Update(); err != nil {
		return errors.BadError("修改失败")
	}
	return nil
}
