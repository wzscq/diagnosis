package report

import (
	"github.com/xuri/excelize/v2"
)

type carInfoItem struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type otherItem struct {
	FurtherCheckedSignalsAnalysis string `json:"FurtherCheckedSignalsAnalysis"`
	PossibleCauses  string `json:"PossibleCauses"`
	RecommendedRecovery string `json:"RecommendedRecovery"`
}

type ecuRecord struct {
	Abnormal string `json:"Abnormal"`
	DtcDescription string `json:"DtcDescription"`
	DtcId string `json:"DtcId"`
	DtcId_State string `json:"DtcId_State"`
	Ecu string `json:"Ecu"`
	Mileage string `json:"Mileage"`
	OtherInfo []otherItem `json:"OtherInfo"`
	Time string  `json:"Time"`
}

type logisticsItem struct {
	Name string `json:"Name"`
	Value string  `json:"value"`
}

type logistics struct {
	EcuName string `json:"EcuName"`
	NameList []logisticsItem `json:"NameList"`
}

type analysisItem struct {
	Item ecuRecord `json:"item"`
	SignalChart map[string]string `json:"signalChart"` 
}

type ecuItem struct {
	Name string `json:"name"`
	Items []ecuRecord  `json:"items"`
	Logistics logistics `json:"logistics"`
}

type ReportContent struct {
	CarInfo  []carInfoItem `json:"carInfo"`
	EcuList map[string]ecuItem `json:"ecuList"`
	AnalysisItems map[string]analysisItem `json:"analysisItems"`
	FileName string `json:"fileName"`
}

func (repo *ReportContent)getReport()(*excelize.File){
	f := excelize.NewFile()
	//Bordered style
	borderStyleNone:=[]excelize.Border{
		{Type: "left", Color: "000000", Style: 0},
		{Type: "top", Color: "000000", Style: 0},
		{Type: "bottom", Color: "000000", Style: 0},
		{Type: "right", Color: "000000", Style: 0},
	}

	borderWhite:=[]excelize.Border{
		{Type: "left", Color: "FFFFFF", Style: 1},
		{Type: "top", Color: "FFFFFF", Style: 1},
		{Type: "bottom", Color: "FFFFFF", Style: 1},
		{Type: "right", Color: "FFFFFF", Style: 1},
	}

	borderStyle:=[]excelize.Border{
		{Type: "left", Color: "888888", Style: 1},
		{Type: "top", Color: "888888", Style: 1},
		{Type: "bottom", Color: "888888", Style: 1},
		{Type: "right", Color: "888888", Style: 1},
	}

	borderSubStyle:=[]excelize.Border{
		{Type: "left", Color: "888888", Style: 1},
		{Type: "top", Color: "888888", Style: 1},
		{Type: "bottom", Color: "888888", Style: 1},
		{Type: "right", Color: "888888", Style: 1},
	}

	//Create style
	styleMainTitle, _ := f.NewStyle(&excelize.Style{
		Border: borderWhite,
		Alignment:&excelize.Alignment{
			Horizontal:"center",
			Vertical:"bottom",
		},
		Font:&excelize.Font {
			Size:14,
			Bold:true,
		},
	})

	styleSubMainTitle, _ := f.NewStyle(&excelize.Style{
		Border: borderWhite,
		Alignment:&excelize.Alignment{
			Horizontal:"center",
			Vertical:"top",
		},
		Font:&excelize.Font {
			Size:12,
			Bold:false,
		},
	})

	styleTitle, _ := f.NewStyle(&excelize.Style{
		Border: borderStyleNone,
		Alignment:&excelize.Alignment{
			Horizontal:"left",
			Vertical:"center",
		},
		Font:&excelize.Font {
			Size:12,
			Bold:true,
		},
	})

	styleSubTitle, _ := f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Alignment:&excelize.Alignment{
			Horizontal:"left",
			Vertical:"center",
			Indent:1,
		},
		Font:&excelize.Font {
			Size:12,
			Bold:true,
		},
		Fill:excelize.Fill{
			Pattern:1,
			Color:[]string{"eeeeee",},
			Type:"pattern",
		},
	})

	styleLabel, _ := f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Alignment:&excelize.Alignment{
			Horizontal:"right",
			Vertical:"center",
			Indent:1,
		},
		Font:&excelize.Font {
			Size:12,
			Bold:false,
		},
		Fill:excelize.Fill{
			Pattern:1,
			Color:[]string{"eeeeee",},
			Type:"pattern",
		},
	})

	styleSubLabel,_:=f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Alignment:&excelize.Alignment{
			Horizontal:"right",
			Vertical:"top",
			Indent:1,
		},
		Font:&excelize.Font {
			Size:12,
			Bold:false,
		},
		Fill:excelize.Fill{
			Pattern:1,
			Color:[]string{"eeeeee",},
			Type:"pattern",
		},
	})

	styleNormal, _ := f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Alignment:&excelize.Alignment{
			Horizontal:"left",
			Vertical:"center",
			Indent:1,
		},
		Font:&excelize.Font {
			Size:12,
			Bold:false,
		},
	})

	styleSubNormal, _ := f.NewStyle(&excelize.Style{
		Border: borderSubStyle,
		Alignment:&excelize.Alignment{
			Horizontal:"left",
			Vertical:"top",
			WrapText:true,
			ShrinkToFit:true,
		},
		Font:&excelize.Font {
			Size:12,
			Bold:false,
		},
	})

	styleYellow, _ := f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Alignment:&excelize.Alignment{
			Horizontal:"left",
			Vertical:"center",
			Indent:1,
		},
		Font:&excelize.Font {
			Size:12,
			Bold:false,
			Color:"FFFF00",
		},
	})

	styleRed, _ := f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Alignment:&excelize.Alignment{
			Horizontal:"left",
			Vertical:"center",
			Indent:1,
		},
		Font:&excelize.Font {
			Size:12,
			Bold:false,
			Color:"FF0000",
		},
	})

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Border: borderStyle,
		Alignment:&excelize.Alignment{
			Horizontal:"center",
			Vertical:"center",
		},
		Font:&excelize.Font {
			Size:12,
			Bold:false,
		},
	})

	styleSheet,_:=f.NewStyle(&excelize.Style{
		Border: borderStyleNone,
		Fill:excelize.Fill{
			Pattern:1,
			Color:[]string{"FFFFFF",},
			Type:"pattern",
		},
	})
	sheetName:="诊断报告"
	// Create a new sheet.
	f.DeleteSheet("Sheet1")
    index := f.NewSheet(sheetName)
    f.SetActiveSheet(index)
	f.SetCellStyle(sheetName,"A1","V10000",styleSheet)

	//8列设置列宽默认100
	f.SetColWidth(sheetName, "A", "H", 20)
	/*for col := 1; col <= 8; col++ {
		cell,_:=excelize.CoordinatesToCellName(col, 1)

	}*/

	// Set value of a cell.
	row:=1
	cellStart,_:=excelize.CoordinatesToCellName(1, row)
	cellEnd,_:=excelize.CoordinatesToCellName(8, row)
	f.MergeCell(sheetName,cellStart,cellEnd)
	f.SetRowHeight(sheetName,row,25)
	f.SetCellStr(sheetName,cellStart,sheetName)	
	f.SetCellStyle(sheetName,cellStart,cellEnd,styleMainTitle)
	row++

	cellStart,_=excelize.CoordinatesToCellName(1, row)
	cellEnd,_=excelize.CoordinatesToCellName(8, row)
	f.MergeCell(sheetName,cellStart,cellEnd)
	f.SetRowHeight(sheetName,row,25)
	f.SetCellStr(sheetName,cellStart,repo.FileName)	
	f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubMainTitle)
	row++

	cellStart,_=excelize.CoordinatesToCellName(1, row)
	cellEnd,_=excelize.CoordinatesToCellName(8, row)
	f.MergeCell(sheetName,cellStart,cellEnd)
	f.SetRowHeight(sheetName,row,25)
	f.SetCellStr(sheetName,cellStart,"1、车辆信息")	
	f.SetCellStyle(sheetName,cellStart,cellEnd,styleTitle)
	row++
	
	for _,carInfoItem :=range repo.CarInfo {
		cellStart,_:=excelize.CoordinatesToCellName(1, row)
		cellEnd,_:=excelize.CoordinatesToCellName(3, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
    	f.SetCellStr(sheetName, cellStart, carInfoItem.Title)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleLabel)
		
		cellStart,_=excelize.CoordinatesToCellName(4, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, carInfoItem.Value)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)
		row++
	}

	cellStart,_=excelize.CoordinatesToCellName(1, row)
	cellEnd,_=excelize.CoordinatesToCellName(8, row)
	f.MergeCell(sheetName,cellStart,cellEnd)
	f.SetRowHeight(sheetName,row,25)
	f.SetCellStr(sheetName,cellStart,"2、故障码概览")		
	f.SetCellStyle(sheetName,cellStart,cellEnd,styleTitle)
	row++

	for _,ecu :=range repo.EcuList {
		cellStart,_:=excelize.CoordinatesToCellName(1, row)
		cellEnd,_:=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName,cellStart,ecu.Name)		
		f.SetRowHeight(sheetName,row,25)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubTitle)
		row++

		cellStart,_=excelize.CoordinatesToCellName(1, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName,cellStart,"物流信息")		
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)
		row++

		for _,lgsItem :=range ecu.Logistics.NameList {
			cellStart,_:=excelize.CoordinatesToCellName(1, row)
			cellEnd,_:=excelize.CoordinatesToCellName(3, row)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, lgsItem.Name)
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleLabel)
			
			cellStart,_=excelize.CoordinatesToCellName(4, row)
			cellEnd,_=excelize.CoordinatesToCellName(8, row)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, lgsItem.Value)
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)
			row++
		}
		
		cellStart,_=excelize.CoordinatesToCellName(1, row)
		f.SetCellStr(sheetName, cellStart, "序号")
		f.SetCellStyle(sheetName,cellStart,cellStart,styleHeader)

		cellStart,_=excelize.CoordinatesToCellName(2, row)
		f.SetCellStr(sheetName, cellStart, "故障代码")
		f.SetCellStyle(sheetName,cellStart,cellStart,styleHeader)

		cellStart,_=excelize.CoordinatesToCellName(3, row)
		cellEnd,_=excelize.CoordinatesToCellName(6, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, "故障内容")
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleHeader)

		cellStart,_=excelize.CoordinatesToCellName(7, row)
		f.SetCellStr(sheetName, cellStart, "故障时刻")
		f.SetCellStyle(sheetName,cellStart,cellStart,styleHeader)

		cellStart,_=excelize.CoordinatesToCellName(8, row)
		f.SetCellStr(sheetName, cellStart, "车辆里程")
		f.SetCellStyle(sheetName,cellStart,cellStart,styleHeader)
		row++

		for index,ecuRec :=range ecu.Items {
			cellStart,_=excelize.CoordinatesToCellName(1, row)
			f.SetCellInt(sheetName, cellStart, index+1)
			f.SetCellStyle(sheetName,cellStart,cellStart,styleNormal)

			cellStart,_=excelize.CoordinatesToCellName(2, row)
			f.SetCellStr(sheetName, cellStart, ecuRec.DtcId)
			if ecuRec.DtcId_State=="0" {
				f.SetCellStyle(sheetName,cellStart,cellStart,styleYellow)
			} else {
				f.SetCellStyle(sheetName,cellStart,cellStart,styleRed)
			}

			cellStart,_=excelize.CoordinatesToCellName(3, row)
			cellEnd,_=excelize.CoordinatesToCellName(6, row)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, ecuRec.DtcDescription)
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)

			cellStart,_=excelize.CoordinatesToCellName(7, row)
			f.SetCellStr(sheetName, cellStart, ecuRec.Time)
			f.SetCellStyle(sheetName,cellStart,cellStart,styleNormal)

			cellStart,_=excelize.CoordinatesToCellName(8, row)
			f.SetCellStr(sheetName, cellStart, ecuRec.Mileage)
			f.SetCellStyle(sheetName,cellStart,cellStart,styleNormal)
			row++	
		}		
	}

	cellStart,_=excelize.CoordinatesToCellName(1, row)
	cellEnd,_=excelize.CoordinatesToCellName(8, row)
	f.MergeCell(sheetName,cellStart,cellEnd)
	f.SetRowHeight(sheetName,row,25)
	f.SetCellStr(sheetName,cellStart,"3、DTC解析")		
	f.SetCellStyle(sheetName,cellStart,cellEnd,styleTitle)
	row++
	
	for _,alyItem:=range repo.AnalysisItems {
		cellStart,_=excelize.CoordinatesToCellName(1, row)
		cellEnd,_=excelize.CoordinatesToCellName(1, row)
		//f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, "故障代码")
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleLabel)

		cellStart,_=excelize.CoordinatesToCellName(2, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, alyItem.Item.DtcId)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)

		row++	

		cellStart,_=excelize.CoordinatesToCellName(1, row)
		cellEnd,_=excelize.CoordinatesToCellName(1, row)
		//f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, "控制器")
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleLabel)

		cellStart,_=excelize.CoordinatesToCellName(2, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, alyItem.Item.Ecu)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)

		row++
		
		cellStart,_=excelize.CoordinatesToCellName(1, row)
		cellEnd,_=excelize.CoordinatesToCellName(1, row)
		//f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, "故障时刻")
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleLabel)

		cellStart,_=excelize.CoordinatesToCellName(2, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, alyItem.Item.Time)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)

		row++

		cellStart,_=excelize.CoordinatesToCellName(1, row)
		cellEnd,_=excelize.CoordinatesToCellName(1, row)
		//f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, "车辆里程")
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleLabel)

		cellStart,_=excelize.CoordinatesToCellName(2, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, alyItem.Item.Mileage)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)

		row++

		cellStart,_=excelize.CoordinatesToCellName(1, row)
		cellEnd,_=excelize.CoordinatesToCellName(1, row)
		//f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, "故障状态")
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleLabel)

		cellStart,_=excelize.CoordinatesToCellName(2, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, alyItem.Item.DtcId_State)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)

		row++

		cellStart,_=excelize.CoordinatesToCellName(1, row)
		cellEnd,_=excelize.CoordinatesToCellName(1, row)
		//f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, "故障内容")
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleLabel)

		cellStart,_=excelize.CoordinatesToCellName(2, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetCellStr(sheetName, cellStart, alyItem.Item.DtcDescription)
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleNormal)

		row++

		for _,otItem:=range alyItem.Item.OtherInfo {
			cellStart,_=excelize.CoordinatesToCellName(1, row)
			cellEnd,_=excelize.CoordinatesToCellName(1, row)
			//f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, "故障原因")
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubLabel)

			cellStart,_=excelize.CoordinatesToCellName(2, row)
			cellEnd,_=excelize.CoordinatesToCellName(4, row)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, otItem.PossibleCauses)
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubNormal)

			cellStart,_=excelize.CoordinatesToCellName(5, row)
			cellEnd,_=excelize.CoordinatesToCellName(5, row)
			//f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, "修复建议")
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubLabel)

			cellStart,_=excelize.CoordinatesToCellName(6, row)
			cellEnd,_=excelize.CoordinatesToCellName(8, row)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, otItem.RecommendedRecovery)
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubNormal)
			f.SetRowHeight(sheetName,row,150)
			row++
		}

		startRow:=row
		for _,otItem:=range alyItem.Item.OtherInfo {
			
			//这里是chart图，暂时空着
			cellStart,_=excelize.CoordinatesToCellName(2, row)
			cellEnd,_=excelize.CoordinatesToCellName(4, row)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, "")
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubNormal)

			cellStart,_=excelize.CoordinatesToCellName(6, row)
			cellEnd,_=excelize.CoordinatesToCellName(8, row)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, otItem.FurtherCheckedSignalsAnalysis)
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubNormal)
			f.SetRowHeight(sheetName,row,150)
			row++
		}

		if row>startRow {
			cellStart,_=excelize.CoordinatesToCellName(1, startRow)
			cellEnd,_=excelize.CoordinatesToCellName(1, row-1)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, "相关信号")
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubLabel)

			cellStart,_=excelize.CoordinatesToCellName(5, startRow)
			cellEnd,_=excelize.CoordinatesToCellName(5, row-1)
			f.MergeCell(sheetName,cellStart,cellEnd)
			f.SetCellStr(sheetName, cellStart, "信号分析")
			f.SetCellStyle(sheetName,cellStart,cellEnd,styleSubLabel)
		}

		cellStart,_=excelize.CoordinatesToCellName(1, row)
		cellEnd,_=excelize.CoordinatesToCellName(8, row)
		f.MergeCell(sheetName,cellStart,cellEnd)
		f.SetRowHeight(sheetName,row,25)
		f.SetCellStr(sheetName,cellStart,"")		
		f.SetCellStyle(sheetName,cellStart,cellEnd,styleTitle)
		row++
	} 	
	
    return f
}