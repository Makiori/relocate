package model

import (
	"relocate/api"
	"relocate/util/errors"

	"github.com/jinzhu/gorm"
)

// 摇珠结果表
type Result struct {
	Model
	DeclarationID         uint   `json:"declaration_id" gorm:"not null;comment:'申报号-外键'"`
	PublicityStatus       bool   `json:"publicity_status" gorm:"not null;comment:'公示状态'"`
	BuildingNo            string `json:"building_no" gorm:"null;comment:'楼号'"`
	RoomNo                string `json:"room_no" gorm:"null;comment:'房间号'"`
	Operator              string `json:"operator" gorm:"not null;comment:'录入结果人员 管理员姓名'"`
	Peoples               string `json:"peoples" gorm:"not null;comment:'被拆迁人(可能有多人)'"`
	CardNumber            string `json:"card_number" gorm:"null;comment:'被拆迁人身份证号码(可能有多人)'"`
	ContractNo            string `json:"contract_no" gorm:"not null;comment:'合同号'"`
	StagingID             uint   `json:"staging_id" gorm:"null;comment:'分期数ID'"`
	PhoneNumber1          string `json:"phone_number1" gorm:"null;comment:'手机号码1'"`
	PhoneNumber2          string `json:"phone_number2" gorm:"null;comment:'手机号码2'"`
	DeclarationHuxingID   uint   `json:"declaration_huxing_id" gorm:"not null;comment:'申报户型ID'"`
	DeclarationHuxingNo   string `json:"declaration_huxing_no" gorm:"not null;comment:'申报户型'"`
	DeclarationBuildingNo string `json:"declaration_huxing_building_no" gorm:"not null;comment:'申报户型栋号'"`
	DeclarationAreaShow   string `json:"declaration_area_show" gorm:"not null;comment:'申报户型面积显示'"`
	DeclarationArea       string `json:"declaration_area" gorm:"not null;comment:'申报面积㎡'"`
	TimeID                uint   `json:"time_id" gorm:"not null;comment:'时段id-外键'"`
	TimeName              string `json:"time_name" gorm:"not null;comment:'冗余 时段表述'"`
}

func (r Result) TableName() string {
	return "result"
}

func initResult() {
	if !db.HasTable(&Result{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Result{}).Error; err != nil {
			panic(err)
		}
		//创建外键
		db.Model(&Result{}).
			AddForeignKey("declaration_id", "declaration(id)", "RESTRICT", "RESTRICT")
	}
}

func (r *Result) Create() error {
	return db.Create(&r).Error
}

type ResultData struct {
	ID                    int    `json:"id"`
	DeclarationID         uint   `json:"declaration_id"`
	DeclarationHuxingID   uint   `json:"declaration_huxing_id"`
	DeclarationHuxingNo   string `json:"declaration_huxing_no"`
	DeclarationBuildingNo string `json:"declaration_huxing_building_no"`
	DeclarationAreaShow   string `json:"declaration_area_show"`
	DeclarationArea       string `json:"declaration_area"`
	BuildingNo            string `json:"building_no"`
	RoomNo                string `json:"room_no"`
	StagingID             uint   `json:"staging_id"`
	StagingName           string `json:"staging_name"`
	ContractNo            string `json:"contract_no"`
	Peoples               string `json:"peoples"`
	Operator              string `json:"operator"`
	PublicityStatus       bool   `json:"publicity_status"`
	TimeID                uint   `json:"time_id"`
	TimeName              string `json:"time_name"`
}

func GetLikeResults(resultFilterBody api.ResultFilterBody) (data *PaginationQ, err error) {
	q := &PaginationQ{
		PageSize: resultFilterBody.PageSize,
		Page:     resultFilterBody.Page,
		Data:     &[]ResultData{},
	}
	status := make([]bool, 0)
	if resultFilterBody.PublicityStatus == -1 {
		status = append(status, true)
		status = append(status, false)
	} else if resultFilterBody.PublicityStatus == 0 {
		status = append(status, false)
	} else {
		status = append(status, true)
	}
	args := "%" + resultFilterBody.FilterName + "%"
	if resultFilterBody.HuxingId == 0 {
		return q.SearchAll(
			db.Table("result r").Select("r.*,"+
				"s.staging_name",
			).Joins("join staging s on "+
				"s.id = r.staging_id",
			).Where("r.staging_id = ? AND "+
				"r.publicity_status IN (?) AND "+
				"(r.card_number LIKE ? OR r.phone_number1 LIKE ? OR r.phone_number2 LIKE ? OR r.contract_no LIKE ? OR r.peoples LIKE ?) and r.deleted_at is null",
				resultFilterBody.StagingId,
				status,
				args,
				args,
				args,
				args,
				args,
			).Order("r.declaration_huxing_no"),
		)
	} else {
		return q.SearchAll(
			db.Table("result r").Select("r.*,"+
				"s.staging_name",
			).Joins("join staging s on "+
				"s.id = r.staging_id",
			).Where("r.staging_id = ? AND "+
				"r.declaration_huxing_id = ? AND "+
				"r.publicity_status IN (?) AND "+
				"(r.card_number LIKE ? OR r.phone_number1 LIKE ? OR r.phone_number2 LIKE ? OR r.contract_no LIKE ? OR r.peoples LIKE ?) and r.deleted_at is null",
				resultFilterBody.StagingId,
				resultFilterBody.HuxingId,
				status,
				args,
				args,
				args,
				args,
				args,
			).Order("r.declaration_huxing_no"),
		)
	}
}

type ResultByStatusData struct {
	DeclarationHuxingNo string `json:"declaration_huxing_no"`
	DeclarationAreaShow string `json:"declaration_area_show"`
	BuildingNo          string `json:"building_no"`
	RoomNo              string `json:"room_no"`
	StagingName         string `json:"staging_name"`
	ContractNo          string `json:"contract_no"`
	Peoples             string `json:"peoples"`
	DeclarationArea     string `json:"declaration_area"`
}

type UserResultData struct {
	ID                  int    `json:"id"`
	DeclarationID       uint   `json:"declaration_id"`
	DeclarationHuxingNo string `json:"declaration_huxing_no"`
	DeclarationAreaShow string `json:"declaration_area_show"`
	DeclarationArea     string `json:"declaration_area"`
	BuildingNo          string `json:"building_no"`
	RoomNo              string `json:"room_no"`
	StagingID           uint   `json:"staging_id"`
	StagingName         string `json:"staging_name"`
	ContractNo          string `json:"contract_no"`
	Peoples             string `json:"peoples"`
	OldAddress          string `json:"old_address"`
}

func FindResultByDeclarationID(declarationID uint) (*Result, error) {
	result := new(Result)
	if err := db.Model(&Result{}).Where("declaration_id = ? and result.deleted_at is null", declarationID).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("该结果不存在")
		}
		return nil, err
	}
	return result, nil
}

func (r *Result) DeleteResultByID() error {
	if err := db.Model(&Result{}).Where("id = ?", r.ID).Delete(&r).Error; err != nil {
		return err
	}
	return nil
}

func FindInResult(contractNo string) ([]Result, error) {
	var results []Result
	err := db.Model(&Result{}).Where("contract_no = ?", contractNo).Find(&results).Error
	return results, err
}

func FindResultByID(resultID uint) (*Result, error) {
	var result Result
	err := db.Model(&Result{}).Where("id = ?", resultID).Take(&result).Error
	return &result, err
}

func FindInResultByID(resultIDList []uint) ([]Result, error) {
	var results []Result
	err := db.Model(&Result{}).Where("id in (?)", resultIDList).Find(&results).Error
	return results, err
}

func (r *Result) FindInResultByDeclarationID(id []int) (*Result, error) {
	err := db.Where("declaration_id in (?)", id).Find(&r).Error
	if err != nil {
		return nil, err
	}
	return r, err
}

func FindResultsByStatus(page, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]Result{},
	}
	return q.SearchAll(db.Model(&Result{}).Where("publicity_status = ?", 1).Order("created_at desc"))
}

func (r *Result) DeleteResultByDeclarationID() error {
	if err := db.Model(&Result{}).Where("declaration_id = ?", r.DeclarationID).Delete(&r).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePublicityResult(d []int, b bool) error {
	sql := db.Model(&Result{}).Where("declaration_id IN (?)", d).Update("publicity_status", b).Error
	//rowsAffected := sql.RowsAffected
	//logging.Infof("更新影响的记录数%d", rowsAffected)
	//logging.Infoln(sql.Error)
	//return sql.Error
	return sql
}
