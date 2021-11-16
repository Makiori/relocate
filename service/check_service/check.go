package check_service

import (
	"fmt"
	"relocate/model"
	"relocate/util/random"
	"relocate/util/times"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func ExportExcel(checkDataList []model.Check) (*excelize.File, string, error) {
	file, err := excelize.OpenFile("./excel/accounting.xlsx")
	if err != nil {
		return nil, "", err
	}
	sheet := "668355"
	for i, check := range checkDataList {
		file.SetCellInt(sheet, fmt.Sprintf("A%d", i+3), i+1)
		file.SetSheetRow(sheet, fmt.Sprintf("B%d", i+3), &[]interface{}{
			check.ContractNo,
			check.SocialCategory,
			check.Peoples,
			check.HouseNumber,
			check.Desc,
			check.TargetPlacementArea,
			check.PlacementOfNonTargetArea,
			check.TemporaryRelocationArea,
			check.TemporaryRelocationSubPlacementOfNonTargetArea,
			check.InitialHQArea,
			check.NonIndexAreaRatio,
			check.IndexAreaRatio,
			check.TemporaryRelocationAreaRatioNonIndex,
		})
		switch check.Rounds {
		case 1:
			file.SetCellValue(sheet, fmt.Sprintf("O%d", i+3), check.BuildingNo+check.RoomNo)
		case 2:
			file.SetCellValue(sheet, fmt.Sprintf("P%d", i+3), check.BuildingNo+check.RoomNo)
		case 3:
			file.SetCellValue(sheet, fmt.Sprintf("R%d", i+3), check.BuildingNo+check.RoomNo)
			switch check.MeasuredFloorArea {
			case 80.4:
				file.SetCellValue(sheet, fmt.Sprintf("AB%d", i+3), check.MeasuredFloorArea)
			case 81.2:
				file.SetCellValue(sheet, fmt.Sprintf("AC%d", i+3), check.MeasuredFloorArea)
			case 100:
				file.SetCellValue(sheet, fmt.Sprintf("AD%d", i+3), check.MeasuredFloorArea)
			}
		}
		if check.Rounds == 1 || check.Rounds == 2 {
			switch check.MeasuredFloorArea {
			case 80.4:
				file.SetCellValue(sheet, fmt.Sprintf("T%d", i+3), check.MeasuredFloorArea)
			case 81.2:
				file.SetCellValue(sheet, fmt.Sprintf("U%d", i+3), check.MeasuredFloorArea)
			case 100:
				file.SetCellValue(sheet, fmt.Sprintf("V%d", i+3), check.MeasuredFloorArea)
			case 122.5:
				file.SetCellValue(sheet, fmt.Sprintf("W%d", i+3), check.MeasuredFloorArea)
			case 122.9:
				file.SetCellValue(sheet, fmt.Sprintf("X%d", i+3), check.MeasuredFloorArea)
			case 139.9:
				file.SetCellValue(sheet, fmt.Sprintf("Y%d", i+3), check.MeasuredFloorArea)
			case 160.1:
				file.SetCellValue(sheet, fmt.Sprintf("Z%d", i+3), check.MeasuredFloorArea)
			case 182.3:
				file.SetCellValue(sheet, fmt.Sprintf("AA%d", i+3), check.MeasuredFloorArea)
			}
		}
		file.SetCellValue(sheet, fmt.Sprintf("S%d", i+3), check.RemainingInitialHQArea)

		file.SetSheetRow(sheet, fmt.Sprintf("AE%d", i+3), &[]interface{}{
			check.MeasuredFloorArea,
			check.UseTargetPlacementArea,
			check.UsePlacementOfNonTargetArea,
			check.UseTemporaryRelocationArea,
			"",
			check.RemainingPlacementOfNonTargetArea,
			check.RemainingTargetPlacementArea,
			check.RemainingTemporaryRelocationArea,
			check.RemainingInitialHQArea,
			check.AmountOfUsedArea,
		})
	}
	fileName := fmt.Sprintf("核算表-%s-%s.xlsx", times.ToStr(), random.String(6))
	return file, fileName, nil
}

func AddCheck(resultID uint, realityArea float64) error {
	return nil
}
