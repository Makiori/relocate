package model

// 管理员日志表
type Logging struct {
	Model
	Username    string `json:"username" gorm:"not null;comment:'操作人'"`
	StagingName string `json:"staging_name" gorm:"not null;comment:'操作期数'"`
	Operation   string `json:"operation" gorm:"not null;comment:'操作'"`
	Details     string `json:"details" gorm:"not null;comment:'详情'"`
}

func (l Logging) TableName() string {
	return "logging"
}

func initLogging() {
	if !db.HasTable(&Logging{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Logging{}).Error; err != nil {
			panic(err)
		}
	}
}

func (l *Logging) Create() error {
	return db.Create(&l).Error
}

func GetLogging() (data *PaginationQ, err error) {
	q := PaginationQ{
		Data: &[]Logging{},
	}
	return q.SearchAll(db.Model(&Logging{}))
}
