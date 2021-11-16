package result_service

import (
	"fmt"
	"relocate/model"
	"relocate/util/errors"
	"relocate/util/random"
	"relocate/util/times"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jinzhu/gorm"
)

func ExportExcel(resultDataList []model.ResultData) (*excelize.File, string, error) {
	file := excelize.NewFile()
	sheet := "Sheet1"
	file.SetCellStr(sheet, "A1", "序号")
	file.SetSheetRow(sheet, "B1", &[]interface{}{
		"结果ID",
		"分期期数ID",
		"分期期数名称",
		"合同号",
		"被拆迁人",
		"申报ID",
		"申报户型ID",
		"申报户型的栋号",
		"申报户型的户型",
		"申报面积㎡",
		"申报面积描述",
		"楼号",
		"房间号",
		"公示状态",
		"录入结果人员",
	})
	for i, resultData := range resultDataList {
		file.SetCellInt(sheet, fmt.Sprintf("A%d", i+2), i+1)
		file.SetSheetRow(sheet, fmt.Sprintf("B%d", i+2), &[]interface{}{
			resultData.ID,
			resultData.StagingID,
			resultData.StagingName,
			resultData.ContractNo,
			resultData.Peoples,
			resultData.DeclarationID,
			resultData.DeclarationHuxingID,
			resultData.DeclarationBuildingNo,
			resultData.DeclarationHuxingNo,
			resultData.DeclarationArea,
			resultData.DeclarationAreaShow,
			resultData.BuildingNo,
			resultData.RoomNo,
			resultData.PublicityStatus,
			resultData.Operator,
		})
	}
	fileName := fmt.Sprintf("摇珠表-%s-%s.xlsx", times.ToStr(), random.String(6))
	//path := fmt.Sprintf("./excel/export/%s", fileName)
	//err := file.SaveAs(path)
	//if err != nil {
	//	return nil, "", err
	//}
	return file, fileName, nil
}

func getDeclartionID(id []int) (*model.Result, error) {
	var r model.Result
	result, err := r.FindInResultByDeclarationID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("有申报号不存在")
		}
		return nil, err
	}
	return result, nil
}

func GetResultByID(id uint) (*model.Result, error) {
	result, err := model.FindResultByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("结果不存在")
		}
		return nil, err
	}
	return result, nil
}

func UpdatePublicityStat(declarationID []int, publicityStatus bool) error {
	if _, err := getDeclartionID(declarationID); err != nil {
		return err
	}
	if err := model.UpdatePublicityResult(declarationID, publicityStatus); err != nil {
		return errors.BadError("（批量）修改状态失败")
	}
	return nil
}
