package model

// 分期，如：西华一期
type Staging struct {
	Model
	StagingName string `json:"staging_name" gorm:"not null;comment:'分期名称'"`
}

func (s Staging) TableName() string {
	return "staging"
}

func initStaging() {
	if !db.HasTable(&Staging{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Staging{}).Error; err != nil {
			panic(err)
		}
	}
}

func FindStagingById(id uint) (*Staging, error) {
	var staging Staging
	if err := db.Where("id = ?", id).First(&staging).Error; err != nil {
		return nil, err
	}
	return &staging, nil
}

type ContractCount struct {
	Id          uint   `json:"id"`
	StagingName string `json:"staging_name"`
	Count       int    `json:"count"`
}

func GetStagingContractCount(page, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]ContractCount{},
	}
	return q.SearchAll(db.Table("staging s").Select("s.*,(select sum(1) from contract c where c.staging_id = s.id) as count"))
}

func GetStaging(page, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]Staging{},
	}
	return q.SearchAll(db.Model(&Staging{}).Order("created_at desc"))
}

func (s *Staging) CreateStaging() error {
	return db.Create(&s).Error
}
