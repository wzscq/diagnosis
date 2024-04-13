package report

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
)

var reportFields = []map[string]interface{}{
	{
		"field": "vehicle_management_code",
	},
	{	
		"field": "vin",
	},
	{	
		"field": "project_num",
	},
	{	
		"field": "specifications",
	},
	{	
		"field": "device_number",
	},
	{	
		"field": "time",
	},
}

func CloseReports(crvClient *crv.CRVClient,token string, reports []string, status map[string]interface{})(int){
	log.Println("close report ...")
	for _, reportID := range reports {
		errorCode := CloseReport(crvClient, token, reportID, status)
		if errorCode != common.ResultSuccess {
			return errorCode
		}
	}
	return common.ResultSuccess
}

func CloseReport(crvClient *crv.CRVClient, token string, reportID string, status map[string]interface{})(int){
	log.Println("close report ...")
	//get report
	commonRep:=crv.CommonReq{
		ModelID:"diag_result",
		Filter:&map[string]interface{}{
			"id":map[string]interface{}{
				"Op.eq":reportID,
			},
		},
		Fields:&reportFields,
	}

	rsp,errorCode:=crvClient.Query(&commonRep,token)

	if errorCode!=common.ResultSuccess {
		return errorCode
	}

	list,_:=rsp.Result["list"].([]interface{})
	if len(list)>0 {
		row,_:=list[0].(map[string]interface{})
		vehicle_management_code:=row["vehicle_management_code"]
		vehicle_management_code_map:=map[string]interface{}{
			"Op.eq":vehicle_management_code,
		}
		if vehicle_management_code==nil {
			vehicle_management_code_map=map[string]interface{}{
				"Op.is":vehicle_management_code,
			}
		}

		project_num:=row["project_num"]
		project_num_map:=map[string]interface{}{
			"Op.eq":project_num,
		}
		if project_num==nil {
			project_num_map=map[string]interface{}{
				"Op.is":project_num,
			}
		}
		specifications:=row["specifications"]
		specifications_map:=map[string]interface{}{
			"Op.eq":specifications,
		}
		if specifications==nil {
			specifications_map=map[string]interface{}{
				"Op.is":specifications,
			}
		}

		device_number:=row["device_number"]
		device_number_map:=map[string]interface{}{
			"Op.eq":device_number,
		}
		if device_number==nil {
			device_number_map=map[string]interface{}{
				"Op.is":device_number,
			}
		}

		time:=row["time"]
		time_map:=map[string]interface{}{
			"Op.eq":time,
		}
		if time==nil {
			time_map=map[string]interface{}{
				"Op.is":time,
			}
		}

		commonRep:=crv.CommonReq{
			ModelID:"diag_result",
			Filter:&map[string]interface{}{
				"vehicle_management_code":vehicle_management_code_map,
				"project_num":project_num_map,
				"specifications":specifications_map,
				"device_number":device_number_map,
				"time":time_map,
			},
			Fields:&[]map[string]interface{}{
				{
					"field":"id",
				},
			},
		}
		rsp,errorCode=crvClient.Query(&commonRep,token)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}

		list,_:=rsp.Result["list"].([]interface{})
		if len(list)>0 {
			var rows []string
			for _,row:=range(list) {
				rowid:=row.(map[string]interface{})["id"].(string)
				rows=append(rows,rowid)
			}

			commonRep:=crv.CommonReq{
				ModelID:"diag_result",
				SelectedRowKeys:&rows,
				List:&[]map[string]interface{}{
					status,
				},
			}
			_,errorCode=crvClient.Update(&commonRep,token)
			return errorCode
		}
	}

	return common.ResultSuccess
}