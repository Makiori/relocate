package model

import "github.com/jinzhu/gorm"

// 时段表
type Time struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null;comment:'时段名称'"`
	OptionalNum uint   `json:"optional_num" gorm:"not null;comment:'可选数'"`
	SelectedNum uint   `json:"selected_num" gorm:"not null;comment:'已选数'"`
}

var db *gorm.DB

func (t Time) TableName() string {
	return "time"
}

func initTime() {
	if !db.HasTable(&Time{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Time{}).Error; err != nil {
			panic(err)
		}
	}
}

func (t *Time) Create() error {
	return db.Create(&t).Error
}

func (t *Time) FindTimeByID(id uint) (*Time, error) {

	err := db.Where("ID = ?", id).First(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func FindAll() ([]Time, error) {
	var times []Time
	err := db.Find(&times).Error
	if err != nil {
		return nil, err
	}
	return times, nil
}

func (t *Time) Delete() error {
	return db.Delete(&t).Error
}

func (t *Time) Update() error {
	return db.Model(&t).Updates(Time{
		Name:        t.Name,
		OptionalNum: t.OptionalNum,
		SelectedNum: t.SelectedNum,
	}).Error
}
