package time_service

import (
	"github.com/jinzhu/gorm"
	"relocate/model"
	"relocate/util/errors"
)

//查找时段
func getTime(id uint) (*model.Time, error) {
	var t model.Time
	time, err := t.FindTimeByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("时段不存在")
		}
		return nil, err
	}
	return time, nil
}
