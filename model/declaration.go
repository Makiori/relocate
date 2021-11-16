package model

import (
	"relocate/api"
	"relocate/util/errors"
	"relocate/util/logging"

	"github.com/jinzhu/gorm"
)

type DeclarationStatus int

const (
	DeclarationOngoing DeclarationStatus = iota
	DeclarationConfirmed
)

func (d DeclarationStatus) String() string {
	switch d {
	case DeclarationOngoing:
		return "进行中"
	case DeclarationConfirmed:
		return "已确认"
	default:
		return "unknown"
	}
}

type WinningStatus int

const (
	WinningNo WinningStatus = iota
	WinningYes
)

func (w WinningStatus) String() string {
	switch w {
	case WinningNo:
		return "不中"
	case WinningYes:
		return "中签"
	default:
		return "unknown"
	}
}

//type HuxingData struct {
//	HuxingNo string `json:"huxing_no"`
//	Area     string `json:"area"`
//	Quantity uint   `json:"quantity"`
//}
//
//type DeclarationHuxingData struct {
//	HuxingData []HuxingData `json:"huxing_data"`
//}

// 申报表
type Declaration struct {
	Model
	//DeclarationHuxingData string            `json:"declaration_huxing_data" gorm:"not null;comment:'申报户型数据-存储json格式文本'"`
	TimeID                uint              `json:"time_id" gorm:"not null;comment:'时段id-外键'"`
	TimeName              string            `json:"time_name" gorm:"not null;comment:'冗余 时段表述'"`
	StagingID             uint              `json:"staging_id" gorm:"not null;comment:'分期数ID-外键'"`
	Rounds                uint              `json:"rounds" gorm:"null;comment:'轮数'"`
	ContractNo            string            `json:"contract_no" gorm:"not null;comment:'合同号-外键'"`
	DeclarationStatus     DeclarationStatus `json:"declaration_status" gorm:"not null;comment:'申报状态'"`
	ActiveState           bool              `json:"active_state" gorm:"not null;comment:'有效状态-是否作废'"`
	WinningStatus         WinningStatus     `json:"winning_status" gorm:"null;comment:'中签状态'"`
	DeclarationHuxingID   uint              `json:"declaration_huxing_id" gorm:"not null;comment:'申报户型ID'"`
	DeclarationHuxingNo   string            `json:"declaration_huxing_no" gorm:"not null;comment:'申报户型'"`
	DeclarationBuildingNo string            `json:"declaration_huxing_building_no" gorm:"not null;comment:'申报户型栋号'"`
	DeclarationAreaShow   string            `json:"declaration_area_show" gorm:"not null;comment:'申报户型面积显示'"`
	DeclarationArea       string            `json:"declaration_area" gorm:"not null;comment:'申报面积㎡'"`
	Trustee               string            `json:"trustee" gorm:"null;comment:'受托人'"`
	TrusteeCardNumber     string            `json:"trustee_card_number" gorm:"null;comment:'受托人身份证号码'"`
	TrusteePhoneNumber    string            `json:"trustee_phone_number" gorm:"null;comment:'受托人手机号码'"`
	TrusteeRelationship   string            `json:"trustee_relationship" gorm:"null;comment:'受托人关系'"`
	Operator              string            `json:"operator" gorm:"not null;comment:'操作人 管理员姓名、登录人姓名'"`
	Printer               string            `json:"printer" gorm:"null;comment:'申报表打印管理员姓名'"`
}

func (d Declaration) TableName() string {
	return "declaration"
}

func initDeclaration() {
	if !db.HasTable(&Declaration{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Declaration{}).Error; err != nil {
			panic(err)
		}
		//创建外键
		db.Model(&Declaration{}).
			AddForeignKey("time_id", "time(id)", "RESTRICT", "RESTRICT")
		db.Model(&Declaration{}).
			AddForeignKey("staging_id", "staging(id)", "RESTRICT", "RESTRICT")
		db.Model(&Declaration{}).
			AddForeignKey("contract_no", "contract(contract_no)", "RESTRICT", "RESTRICT")
		db.Model(&Declaration{}).
			AddForeignKey("declaration_huxing_id", "huxing(id)", "RESTRICT", "RESTRICT")
	}
}

func (d *Declaration) Create() error {
	return db.Create(&d).Error
}

type DeclarationData struct {
	ID                    int               `json:"id"`
	TimeID                uint              `json:"time_id"`
	TimeName              string            `json:"time_name"`
	ContractNo            string            `json:"contract_no"`
	StagingID             uint              `json:"staging_id"`
	StagingName           string            `json:"staging_name"`
	Rounds                uint              `json:"rounds"`
	DeclarationArea       string            `json:"declaration_area"`
	DeclarationHuxingID   uint              `json:"declaration_huxing_id"`
	DeclarationHuxingNo   string            `json:"declaration_huxing_no"`
	DeclarationBuildingNo string            `json:"declaration_huxing_building_no"`
	DeclarationAreaShow   string            `json:"declaration_area_show"`
	ActiveState           bool              `json:"active_state"`
	WinningStatus         WinningStatus     `json:"winning_status"`
	DeclarationStatus     DeclarationStatus `json:"declaration_status"`
	Peoples               string            `json:"peoples"`
	OldAddress            string            `json:"old_address"`
	CardNumber            string            `json:"card_number"`
	HouseNumber           string            `json:"house_number"`
	PhoneNumber1          string            `json:"phone_number1"`
	PhoneNumber2          string            `json:"phone_number2"`
	Trustee               string            `json:"trustee"`
	TrusteeCardNumber     string            `json:"trustee_card_number"`
	TrusteePhoneNumber    string            `json:"trustee_phone_number"`
	TrusteeRelationship   string            `json:"trustee_relationship"`
	SocialCategory        string            `json:"social_category"`
}

func GetAllDeclaration() (data *PaginationQ, err error) {
	q := PaginationQ{
		Data: &[]DeclarationData{},
	}
	return q.SearchAll(
		db.Table("declaration").Select("declaration.*," +
			"contract.*," +
			"staging.*",
		).Joins(
			"join contract on contract.contract_no = declaration.contract_no",
		).Joins(
			"join staging on staging.id = declaration.staging_id",
		).Order("declaration.created_at desc"),
	)
}

type Trustee struct {
	Trustee             string `json:"trustee"`
	TrusteeCardNumber   string `json:"trustee_card_number"`
	TrusteePhoneNumber  string `json:"trustee_phone_number"`
	TrusteeRelationship string `json:"trustee_relationship"`
}

type Detail struct {
	ID                    int               `json:"id"`
	Operator              string            `json:"operator"`
	ContractNo            string            `json:"contract_no"`
	Peoples               string            `json:"peoples"`
	CardNumber            string            `json:"card_number"`
	OldAddress            string            `json:"old_address"`
	InitialHQArea         string            `json:"initial_hq_area"`
	RemainingHQArea       string            `json:"remaining_hq_area"`
	DeclarationArea       string            `json:"declaration_area"`
	DeclarationHuxingID   uint              `json:"declaration_huxing_id"`
	DeclarationHuxingNo   string            `json:"declaration_huxing_no"`
	DeclarationBuildingNo string            `json:"declaration_huxing_building_no"`
	DeclarationAreaShow   string            `json:"declaration_area_show"`
	DeclarationStatus     DeclarationStatus `json:"declaration_status"`
	WinningStatus         WinningStatus     `json:"winning_status"`
	ActiveState           bool              `json:"active_state"`
	TimeID                uint              `json:"time_id"`
	TimeName              string            `json:"time_name"`
	SocialCategory        string            `json:"social_category"`
	HouseNumber           string            `json:"house_number"`
	Trustee
}

type ContractInfo struct {
	ContractNo      string `json:"contract_no"`
	Peoples         string `json:"peoples"`
	CardNumber      string `json:"card_number"`
	OldAddress      string `json:"old_address"`
	InitialHQArea   string `json:"initial_hq_area"`
	RemainingHQArea string `json:"remaining_hq_area"`
}

type DeclarationItem struct {
	ID                    int               `json:"id"`
	Operator              string            `json:"operator"`
	Peoples               string            `json:"peoples"`
	DeclarationArea       string            `json:"declaration_area"`
	DeclarationHuxingID   uint              `json:"declaration_huxing_id"`
	DeclarationHuxingNo   string            `json:"declaration_huxing_no"`
	DeclarationBuildingNo string            `json:"declaration_huxing_building_no"`
	DeclarationAreaShow   string            `json:"declaration_area_show"`
	StagingID             uint              `json:"staging_id"`
	StagingName           string            `json:"staging_name"`
	DeclarationStatus     DeclarationStatus `json:"declaration_status"`
	WinningStatus         WinningStatus     `json:"winning_status"`
	ActiveState           bool              `json:"active_state"`
}

func GetActiveStateDeclaration(contractNo string, stagingID, rounds uint) (data []Declaration, err error) {
	if err = db.Model(&Declaration{}).Where("active_state = ? "+
		"and contract_no = ? "+
		"and staging_id = ? "+
		"and rounds = ?",
		true,
		contractNo,
		stagingID, rounds).Find(&data).Error; err != nil {
		return
	}
	return
}

func FindDeclarationByID(declarationID uint) (declaration *Declaration, err error) {
	declaration = new(Declaration)
	err = db.Model(&Declaration{}).Where("id = ?", declarationID).Find(&declaration).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("该申报不存在")
		}
		return nil, err
	}
	return
}

func GetDeclarationListByID(id []uint) ([]*Declaration, error) {
	var data []*Declaration
	if err := db.Model(&Declaration{}).Where("id in (?)", id).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Declaration) Update() error {
	sql := db.Model(&Declaration{}).Save(&d)
	if err := sql.Error; err != nil {
		return err
	}
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return nil
}

func (d *Declaration) UpdateDeclaration() error {
	sql := db.Model(d).Where("id = ?", d.ID).Updates(&d)
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

func (d *Declaration) UpdateDeclarationStatus(id uint, status DeclarationStatus) error {
	sql := db.Model(&Declaration{}).Where("id = ?", id).Update("declaration_status = ?", status)
	if err := sql.Error; err != nil {
		return err
	}
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return nil
}

func (d *Declaration) UpdateWinningStatus(id uint, status WinningStatus) error {
	sql := db.Model(&Declaration{}).Where("id = ?", id).Update("winning_status = ?", status)
	if err := sql.Error; err != nil {
		return err
	}
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return nil
}

func (d *Declaration) DeleteDeclaration() error {
	return db.Delete(&d).Error
}

func GetDeclaration(stagingId uint, huxing_id uint, time_id uint, declaration_status int, winning_status int, active_state int, filterName string, page uint, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]Declaration{},
	}
	arg1 := "%" + filterName + "%"
	return q.SearchAll(db.Model(&Declaration{}).Where("staging_id = ? OR declaration_huxing_id = ? OR time_id = ? OR declaration_status = ? OR winning_status = ? OR active_state = ? AND (trustee LIKE ? OR trustee_card_number LIKE ? OR trustee_phone_number LIKE ?)", stagingId, huxing_id, time_id, declaration_status, winning_status, active_state, arg1, arg1, arg1).Order("created_at desc"))
}

type DeclarationDetail struct {
	DeclarationData
	Detail
}

func GetDeclarationDetail(contractNo string, stagingId string, declarationID int) (data interface{}, err error) {
	var date []*Detail
	if err := db.Model(&Declaration{}).Where("contract_no = ? AND staging_id = ? AND id = ?", contractNo, stagingId, declarationID).
		Scan(&date).Order("created_at desc").Error; err != nil {
		return nil, err
	}
	return date, err
}

func (d *Declaration) FindDeclartionByStagingId(stagingid uint) (*Declaration, error) {
	err := db.Where("staging_id = ?", stagingid).First(&d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}

func GetLikeDeclaration(declarationFilterBody api.DeclarationFilterBody) (data *PaginationQ, err error) {
	q := &PaginationQ{
		PageSize: declarationFilterBody.PageSize,
		Page:     declarationFilterBody.Page,
		Data:     &[]DeclarationData{},
	}
	status := make([]bool, 0)
	if declarationFilterBody.DeclarationStatus == -1 {
		status = append(status, true)
		status = append(status, false)
	} else if declarationFilterBody.DeclarationStatus == 0 {
		status = append(status, true)
	} else {
		status = append(status, false)
	}

	if declarationFilterBody.HuxingId == 0 {
		return q.SearchAll(
			db.Table("declaration d").
				Select("d.*, c.peoples , c.old_address, c.card_number, c.house_number, c.phone_number1, c.phone_number2").Joins("join contract c on c.contract_no = d.contract_no").
				Where("d.staging_id = ? AND d.declaration_status IN (?) AND d.deleted_at is null", declarationFilterBody.StagingId, status).Order("d.declaration_huxing_no"),
		)
	} else {
		return q.SearchAll(
			db.Table("declaration d").
				Select("d.*, c.peoples , c.old_address, c.card_number, c.house_number, c.phone_number1, c.phone_number2").Joins("join contract c on c.contract_no = d.contract_no").
				Where("d.staging_id = ? AND d.declaration_huxing_id = ? AND d.declaration_status IN (?) AND d.deleted_at is null", declarationFilterBody.StagingId, declarationFilterBody.HuxingId, status).Order("d.declaration_huxing_no"),
		)
	}

}

func EnterResult(declarationid uint, buildingno string, roomno string, status int) error {
	sql := db.Model(&Declaration{}).Where("id = ?", declarationid).Update("winning_status", status)
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

func FindAllDeclaration(contractNo string) (data interface{}, err error) {
	var declaration []Declaration
	if err := db.Model(&Declaration{}).Where("contract_no = ?", contractNo).Find(&declaration).Error; err != nil {
		return nil, err
	}
	return declaration, nil

}

func UpdataDeclarationActive(declarationID uint, state bool) error {
	sql := db.Model(&Declaration{}).Where("id = ?", declarationID).Update("active_state", state)
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}
