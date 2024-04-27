package send

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
)

var QueryConfigFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "state_value"},
	{"field": "wait_time"},
	{
		"field": "ecu_state_machine",
		"fieldType": "many2one",
		"relatedModelID": "diag_signal",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "name"},
			{"field": "can_id"},
			{"field": "byte_order"},
			{"field": "start_addr"},
			{"field": "len"},
			{"field": "factor"},
			{"field": "offset"},
			{"field": "pdu_id"},
		},
	},
	{
		"field": "mileage_name",
		"fieldType": "many2one",
		"relatedModelID": "diag_signal",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "name"},
			{"field": "can_id"},
			{"field": "byte_order"},
			{"field": "start_addr"},
			{"field": "len"},
			{"field": "factor"},
			{"field": "offset"},
			{"field": "pdu_id"},
		},
	},
}

type diagConf struct {
	CRVClient *crv.CRVClient
	Conf map[string]interface{}
}

func getDiagConf(crvClient *crv.CRVClient,platform,token string)(*diagConf,int){
	dc:=&diagConf{
		CRVClient:crvClient,
	}

	rsp,errorCode:=dc.queryDiagConfig(platform,token)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	dc.Conf,errorCode=dc.convertDiagConf(rsp)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	return dc,common.ResultSuccess
}

func (dc *diagConf)queryDiagConfig(platform,token string)(*crv.CommonRsp,int){
	log.Println("start queryDiagConfig")
	commonReq:=crv.CommonReq{
		ModelID:"diag_exe_config",
		Fields:&QueryConfigFields,
		Filter:&map[string]interface{}{
			"platform_id":map[string]interface{}{
				"Op.eq":platform,
			},
		},
		Pagination:&crv.Pagination{
			Current:1,
			PageSize:1,
		},
	}

	return dc.CRVClient.Query(&commonReq,token)
}

func (dc *diagConf)convertDiagConf(queryResult *crv.CommonRsp)(map[string]interface{},int){
	log.Println("start convertDiagConf")
	list,ok:=queryResult.Result["list"]
	if !ok {
		log.Println("convertDiagConf queryResult no list")
		return nil,common.ResultNoDiagConf
	}

	listMap,ok:=list.([]interface{})
	if !ok || len(listMap)<=0 {
		log.Println("convertDiagConf queryResult no list")
		return nil,common.ResultNoDiagConf
	}

	row,ok:=listMap[0].(map[string]interface{})
	if !ok {
		log.Println("convertDiagConf convert list[0] to map error")
		return nil,common.ResultNoDiagConf
	}

	ecuMachine,ok:=row["ecu_state_machine"].(map[string]interface{})
	if !ok {
		log.Println("convertDiagConf convert ecu_state_machine to map error")
		return nil,common.ResultNoDiagConf
	}

	conf:=map[string]interface{}{}
	//conf["TriggerCanId"]=row["trigger_can_id"]
	ecuMachineLst,ok:=ecuMachine["list"].([]interface{})
	if ok && len(ecuMachineLst)>0 {
		ecuMachine,_:=ecuMachineLst[0].(map[string]interface{})
		executionPara:=map[string]interface{}{}
		executionPara["StateValue"]=row["state_value"]
		executionPara["WaitTime"]=row["wait_time"]
		executionPara["EcuStateMachine"]=ecuMachine["name"]
		
		executionPara["CanID"]=ecuMachine["can_id"]
		executionPara["ByteOrder"]=ecuMachine["byte_order"]
		executionPara["StartAddr"]=ecuMachine["start_addr"]
		executionPara["Len"]=ecuMachine["len"]
		executionPara["Factor"]=ecuMachine["factor"]
		executionPara["Offset"]=ecuMachine["offset"]
		executionPara["PduID"]=ecuMachine["pdu_id"]
		conf["ExecutionPara"]=executionPara
	}

	mileage,ok:=row["mileage_name"].(map[string]interface{})
	if !ok {
		log.Println("convertDiagConf convert mileage_name to map error")
		return nil,common.ResultNoDiagConf
	}
	mileageLst,ok:=mileage["list"].([]interface{})
	if ok && len(mileageLst)>0 {
		mileage,_:=mileageLst[0].(map[string]interface{})
		mileagePara:=map[string]interface{}{}
		mileagePara["MileageName"]=mileage["name"]

		mileagePara["CanID"]=mileage["can_id"]
		mileagePara["ByteOrder"]=mileage["byte_order"]
		mileagePara["StartAddr"]=mileage["start_addr"]
		mileagePara["Len"]=mileage["len"]
		mileagePara["Factor"]=mileage["factor"]
		mileagePara["Offset"]=mileage["offset"]
		mileagePara["PduID"]=mileage["pdu_id"]
		conf["MileagePara"]=mileagePara
	}

	return conf,common.ResultSuccess
}