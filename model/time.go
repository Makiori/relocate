package model

import (
	"relocate/util/logging"
)

// 时段表
type Time struct {
	Model
	Name        string `json:"name" gorm:"not null;comment:'时段名称'"`
	OptionalNum uint   `json:"optional_num" gorm:"not null;comment:'可选数'"`
	SelectedNum uint   `json:"selected_num" gorm:"not null;comment:'已选数'"`
	StagingId   uint   `json:"staging_id" gorm:"not null;comment:'分期id'"`
}

func (t Time) TableName() string {
	return "time"
}

func initTime() {
	if !db.HasTable(&Time{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Time{}).Error; err != nil {
			panic(err)
		}
		db.Model(&Time{}).
			AddForeignKey("staging", "staging(id)", "RESTRICT", "RESTRICT")
	}
}

func (t *Time) FindTimeByID(id uint) (*Time, error) {

	err := db.Where("ID = ?", id).First(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Time) Update() error {
	sql := db.Model(&t).Updates(Time{
		Name:        t.Name,
		OptionalNum: t.OptionalNum,
		SelectedNum: t.SelectedNum,
	})
	logging.Infof("影响的行数为%d", sql.RowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}
