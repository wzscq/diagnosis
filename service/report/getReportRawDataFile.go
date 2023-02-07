package report

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
)

var queryFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "raw_data"},
}

func getReportRawdataFile(modelID,reportID,token string,crvClinet *crv.CRVClient)(string,int){
	commonRep:=crv.CommonReq{
		ModelID:modelID,
		Filter:&map[string]interface{}{
			"id":reportID,
		},
		Fields:&queryFields,
	}

	rsp,err:=crvClinet.Query(&commonRep,token)
	if err!=common.ResultSuccess {
		return "",err
	}
	list,_:=rsp.Result["list"].([]interface{})
	row,_:=list[0].(map[string]interface{})
	filePath,_:=row["raw_data"].(string)
	return filePath,common.ResultSuccess
}