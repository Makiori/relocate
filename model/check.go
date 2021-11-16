package model

import (
	"relocate/api"
	"relocate/util/times"
)

//核算
type Check struct {
	ID        uint           `json:"id" gorm:"primary_key;comment:'ID'"`
	CreatedAt times.JsonTime `json:"created_at" gorm:"not null;comment:'创建时间'"`
	UpdatedAt times.JsonTime `json:"updated_at" gorm:"comment:'更新时间'"`

	ContractNo              string  `json:"contract_no" gorm:"not null;comment:'合同号'"`
	SocialCategory          string  `json:"social_category" gorm:"null;comment:'社别'"`
	Peoples                 string  `json:"peoples" gorm:"not null;comment:'被拆迁人(可能有多人)'"`
	HouseNumber             string  `json:"house_number" gorm:"null;comment:'房屋栋号'"`
	Desc                    string  `json:"desc" gorm:"null;comment:'备注'"`
	InitialHQArea           float64 `json:"initial_hq_area" gorm:"null;comment:'应得补偿安置总面积'"`
	TargetPlacementArea     float64 `json:"target_placement_area" gorm:"null;comment:'指标安置面积'"`
	TemporaryRelocationArea float64 `json:"temporary_relocation_area" gorm:"null;comment:'计算临迁费面积'"`

	ResultList

	// 计算字段
	PlacementOfNonTargetArea                       float64 `json:"placement_of_non_target_area" gorm:"null;comment:'安置非指标面积'"`
	NonIndexAreaRatio                              float64 `json:"non_index_area_ratio" gorm:"null;comment:'非指标面积占比'"`
	IndexAreaRatio                                 float64 `json:"index_area_ratio" gorm:"null;comment:'指标面积占比'"`
	TemporaryRelocationAreaRatioNonIndex           float64 `json:"temporary_relocation_area_ratio_non_index" gorm:"null;comment:'临迁费面积占安置非指标面积比例'"`
	RemainingResettlementArea                      float64 `json:"remaining_resettlement_area" gorm:"null;comment:'剩余应得回迁安置面积'"`
	TemporaryRelocationSubPlacementOfNonTargetArea float64 `json:"temporary_relocation_sub_placement_of_non_target_area" gorm:"null;comment:'计算临迁费面积与非指标面积之差'"`
	MeasuredFloorArea                              float64 `json:"measured_floor_area" gorm:"null;comment:'中签房号实测建筑面积'"`
	UseTargetPlacementArea                         float64 `json:"use_target_placement_area" gorm:"null;comment:'已使用的安置指标面积'"`
	UsePlacementOfNonTargetArea                    float64 `json:"use_placement_of_non_target_area" gorm:"null;comment:'已使用的安置非指标面积'"`
	UseTemporaryRelocationArea                     float64 `json:"use_temporary_relocation_area" gorm:"null;comment:'已使用的计算临迁安置费面积'"`
	RemainingPlacementOfNonTargetArea              float64 `json:"remaining_placement_of_non_target_area" gorm:"null;comment:'剩余安置非指标面积'"`
	RemainingTargetPlacementArea                   float64 `json:"remaining_target_placement_area" gorm:"null;comment:'剩余安置指标面积'"`
	RemainingTemporaryRelocationArea               float64 `json:"remaining_temporary_relocation_area" gorm:"null;comment:'剩余计算临迁费面积'"`
	RemainingInitialHQArea                         float64 `json:"remaining_initial_hq_area" gorm:"null;comment:'剩余应得回迁安置总面积'"`
	AmountOfUsedArea                               float64 `json:"amount_of_used_area" gorm:"null;comment:'购买已使用指安置标面积的金额1000元/㎡'"`
}

func (c Check) TableName() string {
	return "check"
}

func (c *Check) Create() error {
	return db.Create(&c).Error
}

func (c *Check) Update() error {
	return db.Model(&Check{}).
		Where("result_id = ?", c.ResultID).
		Updates(&c).Error
}

func FindCheckByResultID(resultID uint) (*Check, error) {
	var c Check
	err := db.Model(&Check{}).Where("result_id = ?", resultID).Take(&c).Error
	return &c, err
}

func GetLikeCheck(checkFilterBody api.CheckFilterBody) (data *PaginationQ, err error) {
	c := &PaginationQ{
		PageSize: checkFilterBody.PageSize,
		Page:     checkFilterBody.Page,
		Data:     &[]Check{},
	}
	args := "%" + checkFilterBody.FilterName + "%"
	data, err = c.SearchAll(
		db.Model(&Check{}).Where("contract_no LIKE ? OR peoples LIKE ?", args, args),
	)
	if err != nil {
		return nil, err
	}
	return
}
