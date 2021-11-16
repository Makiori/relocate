package model

import (
	"github.com/jinzhu/gorm"
	"relocate/util/errors"
	"relocate/util/logging"
)

const (
	NowStagingConfig     string = "NowStagingConfig"
	NowRoundsConfig      string = "NowRoundsConfig"
	HuxingOptionalConfig string = "HuxingOptionalConfig"
	IncludeTotalConfig   string = "IncludeTotalConfig" //全局配置摇珠不中是否计入总数，0：不计入总数；1：计入总数
)

// 全局配置表
type Config struct {
	Key   string `json:"key" gorm:"unique;not null;comment:'key'"`
	Value uint   `json:"value" gorm:"not null;comment:'value'"`
}

func (c Config) TableName() string {
	return "config"
}

func initConfig() {
	if !db.HasTable(&Config{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Config{}).Error; err != nil {
			panic(err)
		}
		nowStagingConfig := Config{
			Key:   NowStagingConfig,
			Value: 0,
		}
		nowStagingConfig.Create()
		huxingOptionalConfig := Config{
			Key:   HuxingOptionalConfig,
			Value: 0,
		}
		huxingOptionalConfig.Create()
		nowRoundsConfig := Config{
			Key:   NowRoundsConfig,
			Value: 1,
		}
		nowRoundsConfig.Create()
		includeTotalConfig := Config{
			Key:   IncludeTotalConfig,
			Value: 0,
		}
		includeTotalConfig.Create()
	}
}

func (c *Config) Create() error {
	return db.Create(&c).Error
}

func SettingNowStagingConfig(stagingId uint) error {
	logging.Infoln(stagingId)
	if err := db.Model(&Config{}).Where("`key` = ?", NowStagingConfig).Update("value", stagingId).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func GetNowRoundsConfig() (uint, error) {
	var config Config
	if err := db.Where("`key` = ?", NowRoundsConfig).First(&config).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return config.Value, nil
}

func SettingNowRoundsConfig(rounds uint) error {
	if err := db.Model(&Config{}).Where("`key` = ?", NowRoundsConfig).Update("value", rounds).
		Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func GetNowStagingConfig() (uint, error) {
	var config Config
	if err := db.Where("`key` = ?", NowStagingConfig).First(&config).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return config.Value, nil
}

//配置再不中签的情况下，申报表的有效状态的改变
func SettingIncludeTotalConfig(state uint) error {
	if err := db.Model(&Config{}).Where("`key` = ?", IncludeTotalConfig).Update("value", state).Error; err != nil {
		return err
	}
	return nil
}

//获取当前摇珠不中是否计入申报条数状态
func GetIncludeTotalConfig() (state uint, err error) {
	var config Config
	if err = db.Model(&Config{}).Where("`key` = ?", IncludeTotalConfig).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.BadError("不存在摇珠是否计入总数全局配置")
		}
		return
	}
	state = config.Value
	return
}
