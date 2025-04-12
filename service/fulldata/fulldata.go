package fulldata

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
)

var queryFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "file_name"},
}

func GetFullDataFileName(dataID string, token string, crvClient *crv.CRVClient) (string, int) {
	commonRep:=crv.CommonReq{
		ModelID:"full_data_rec",
		Filter:&map[string]interface{}{
			"id":dataID,
		},
		Fields:&queryFields,
	}

	rsp,errCode:=crvClient.Query(&commonRep,token)
	if errCode!=common.ResultSuccess {
		return "",errCode
	}
	list,_:=rsp.Result["list"].([]interface{})
	row,_:=list[0].(map[string]interface{})
	fileName,_:=row["file_name"].(string)
	return fileName,common.ResultSuccess
}