package check_service

import (
	"fmt"
	"relocate/model"
	"relocate/service/contract_service"
	"relocate/service/declaration_service"
	"relocate/service/huxing_service"
	"relocate/service/result_service"
	"relocate/util/random"
	"relocate/util/times"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
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
	fileName := fmt.Sprintf("?????????-%s-%s.xlsx", times.ToStr(), random.String(6))
	return file, fileName, nil
}

func AddCheck(resultID uint, realityArea float64) error {
	result, err := result_service.GetResultByID(resultID)
	if err != nil {
		return err
	}

	declaration, err := declaration_service.GetDeclarationById(result.DeclarationID)
	if err != nil {
		return err
	}

	contract, err := contract_service.GetContractById(declaration.ContractNo)
	if err != nil {
		return err
	}

	huxing := new(model.Huxing)
	huxing, err = huxing_service.Gethuxing(result.DeclarationHuxingID)

	if err != nil {
		return err
	}
	resultList := model.ResultList{
		ResultID:            result.ID,
		BuildingNo:          result.BuildingNo,
		RoomNo:              result.RoomNo,
		DeclarationHuxingID: huxing.ID,
		Rounds:              huxing.Rounds,
		DeclarationArea:     result.DeclarationArea,
	}

	//??????

	//???????????????????????????
	initialHQAreaDecimal := decimal.NewFromFloat(contract.InitialHQArea)
	//??????????????????
	targetPlacementAreaDecimal := decimal.NewFromFloat(contract.TargetPlacementArea)
	//?????????????????????
	temporaryRelocationAreaDecimal := decimal.NewFromFloat(contract.TemporaryRelocationArea)

	//?????????????????????
	placementOfNonTargetAreaDecimal := initialHQAreaDecimal.Sub(targetPlacementAreaDecimal)
	placementOfNonTargetArea, _ := placementOfNonTargetAreaDecimal.Float64()

	//?????????????????????
	nonIndexAreaRatioDecimal := decimal.NewFromFloat(0)
	if f, _ := initialHQAreaDecimal.Float64(); f != 0 {
		nonIndexAreaRatioDecimal = placementOfNonTargetAreaDecimal.Div(initialHQAreaDecimal)
	}
	nonIndexAreaRatio, _ := nonIndexAreaRatioDecimal.Float64()

	//??????????????????
	indexAreaRatioDecimal := decimal.NewFromFloat(0)
	if f, _ := initialHQAreaDecimal.Float64(); f != 0 {
		indexAreaRatioDecimal = targetPlacementAreaDecimal.Div(initialHQAreaDecimal)
	}
	indexAreaRatio, _ := indexAreaRatioDecimal.Float64()

	//?????????????????????????????????????????????
	temporaryRelocationAreaRatioNonIndexDecimal := decimal.NewFromFloat(0)
	if f, _ := placementOfNonTargetAreaDecimal.Float64(); f != 0 {
		temporaryRelocationAreaRatioNonIndexDecimal = temporaryRelocationAreaDecimal.Div(placementOfNonTargetAreaDecimal)
	}
	temporaryRelocationAreaRatioNonIndex, _ := temporaryRelocationAreaRatioNonIndexDecimal.Float64()

	//??????????????????????????????
	var measuredFloorArea float64
	measuredFloorAreaDecimal := decimal.NewFromFloat(0)
	declarationAreaDecimal, _ := decimal.NewFromString(resultList.DeclarationArea)
	measuredFloorAreaDecimal = measuredFloorAreaDecimal.Add(declarationAreaDecimal)
	if realityArea > 0 {
		measuredFloorArea = realityArea
	} else {
		measuredFloorArea, _ = measuredFloorAreaDecimal.Float64()
	}

	//??????????????????????????????
	remainingResettlementAreaDecimal := initialHQAreaDecimal.Sub(measuredFloorAreaDecimal)
	remainingResettlementArea, _ := remainingResettlementAreaDecimal.Float64()

	//?????????????????????????????????????????????
	temporaryRelocationSubPlacementOfNonTargetArea, _ := temporaryRelocationAreaDecimal.Sub(placementOfNonTargetAreaDecimal).Float64()

	//??????????????????????????????
	useTargetPlacementAreaDecimal := measuredFloorAreaDecimal.Mul(indexAreaRatioDecimal)
	useTargetPlacementArea, _ := useTargetPlacementAreaDecimal.Float64()

	//?????????????????????????????????
	usePlacementOfNonTargetAreaDecimal := measuredFloorAreaDecimal.Mul(nonIndexAreaRatioDecimal)
	usePlacementOfNonTargetArea, _ := usePlacementOfNonTargetAreaDecimal.Float64()

	//???????????????????????????????????????
	useTemporaryRelocationAreaDecimal := usePlacementOfNonTargetAreaDecimal.Mul(temporaryRelocationAreaRatioNonIndexDecimal)
	useTemporaryRelocationArea, _ := useTemporaryRelocationAreaDecimal.Float64()

	//???????????????????????????
	remainingPlacementOfNonTargetAreaDecimal := placementOfNonTargetAreaDecimal.Sub(usePlacementOfNonTargetAreaDecimal)
	remainingPlacementOfNonTargetArea, _ := remainingPlacementOfNonTargetAreaDecimal.Float64()

	//????????????????????????
	remainingTargetPlacementAreaDecimal := targetPlacementAreaDecimal.Sub(useTargetPlacementAreaDecimal)
	remainingTargetPlacementArea, _ := remainingTargetPlacementAreaDecimal.Float64()

	//???????????????????????????
	remainingTemporaryRelocationAreaDecimal := temporaryRelocationAreaDecimal.Sub(useTemporaryRelocationAreaDecimal)
	remainingTemporaryRelocationArea, _ := remainingTemporaryRelocationAreaDecimal.Float64()

	//?????????????????????????????????
	remainingInitialHQAreaDecimal := initialHQAreaDecimal.Sub(measuredFloorAreaDecimal)
	remainingInitialHQArea, _ := remainingInitialHQAreaDecimal.Float64()

	//??????????????????????????????????????????1000???/???
	amountOfUsedAreaDecimal := useTargetPlacementAreaDecimal.Mul(decimal.NewFromFloat(1000))
	amountOfUsedArea, _ := amountOfUsedAreaDecimal.Float64()

	check := model.Check{
		ContractNo:              contract.ContractNo,
		SocialCategory:          contract.SocialCategory,
		Peoples:                 contract.Peoples,
		HouseNumber:             contract.HouseNumber,
		Desc:                    contract.Desc,
		InitialHQArea:           contract.InitialHQArea,
		TargetPlacementArea:     contract.TargetPlacementArea,
		TemporaryRelocationArea: contract.TemporaryRelocationArea,
		ResultList:              resultList,

		PlacementOfNonTargetArea:                       placementOfNonTargetArea,
		NonIndexAreaRatio:                              nonIndexAreaRatio,
		IndexAreaRatio:                                 indexAreaRatio,
		TemporaryRelocationAreaRatioNonIndex:           temporaryRelocationAreaRatioNonIndex,
		RemainingResettlementArea:                      remainingResettlementArea,
		TemporaryRelocationSubPlacementOfNonTargetArea: temporaryRelocationSubPlacementOfNonTargetArea,
		MeasuredFloorArea:                              measuredFloorArea,
		UseTargetPlacementArea:                         useTargetPlacementArea,
		UsePlacementOfNonTargetArea:                    usePlacementOfNonTargetArea,
		UseTemporaryRelocationArea:                     useTemporaryRelocationArea,
		RemainingPlacementOfNonTargetArea:              remainingPlacementOfNonTargetArea,
		RemainingTargetPlacementArea:                   remainingTargetPlacementArea,
		RemainingTemporaryRelocationArea:               remainingTemporaryRelocationArea,
		RemainingInitialHQArea:                         remainingInitialHQArea,
		AmountOfUsedArea:                               amountOfUsedArea,
	}
	_, err = model.FindCheckByResultID(check.ResultID)
	if err == gorm.ErrRecordNotFound {
		_ = check.Create()
	} else {
		_ = check.Update()
	}
	return nil
}
