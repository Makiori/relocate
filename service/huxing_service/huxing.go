package huxing_service

import (
	"relocate/model"
	"relocate/util/errors"

	"github.com/jinzhu/gorm"
)

//查找户型
func gethuxing(id uint) (*model.Huxing, error) {
	var t model.Huxing
	huxing, err := t.FindHuxingByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("户型不存在")
		}
		return nil, err
	}
	return huxing, nil
}

//新增户型
func AddHuxing(area string, areaShow string, buildingNo string, huxingNo string, maximum uint, quantity uint, rounds uint, stagingId uint) error {
	huxing := model.Huxing{
		StagingID:  stagingId,
		Area:       area,
		AreaShow:   areaShow,
		BuildingNo: buildingNo,
		HuxingNo:   huxingNo,
		Maximum:    maximum,
		Quantity:   quantity,
		Rounds:     rounds,
	}
	if err := huxing.CreateHuxing(); err != nil {
		return errors.BadError("新增户型失败")
	}
	return nil
}

//删除户型
func DeleteHuxing(Id uint) error {
	huxing := model.Huxing{
		Model: model.Model{ID: Id},
	}
	if _, err := gethuxing(Id); err != nil {
		return err
	}
	if err := huxing.DeletedHuxing(); err != nil {
		return errors.BadError("删除户型失败")
	}
	return nil
}

//修改户型
func UpdateHuxing(area string, areaShow string, buildingNo string, huxingNo string, maximum uint, quantity uint, rounds uint, id uint) error {
	huxing := model.Huxing{
		Model:      model.Model{ID: id},
		Area:       area,
		AreaShow:   areaShow,
		BuildingNo: buildingNo,
		HuxingNo:   huxingNo,
		Maximum:    maximum,
		Quantity:   quantity,
		Rounds:     rounds,
	}
	if _, err := gethuxing(id); err != nil {
		return err
	}
	if err := huxing.UpdateHuxing(); err != nil {
		return errors.BadError("修改户型失败")
	}
	return nil
}
