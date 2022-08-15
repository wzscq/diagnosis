package send

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
	"fmt"
)

var QueryDTCSignalFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "dtc_id"},
	{
		"field": "signals",
		"fieldType": "many2many",
		"relatedModelID": "diag_signal",
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		 	{"field": "name"},
			{"field": "can_id"},
			{"field": "start_addr"},
			{"field": "pdu_id"},
		},
	},
}

type dtcList struct {
	CRVClient *crv.CRVClient
	Ecus []string
	DtcList []interface{}
}

func getDtcList(ecus []string,crvClient *crv.CRVClient)(*dtcList,int){
	dtc:=&dtcList{
		CRVClient:crvClient,
		Ecus:ecus,
	}

	log.Println("getDtcList",ecus)

	rsp,errorCode:=dtc.queryDtcList()
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	dtc.DtcList,errorCode=dtc.convertDtcList(rsp)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	return dtc,common.ResultSuccess
}

func (dtc *dtcList)queryDtcList()(*crv.CommonRsp,int){
	log.Println("start queryDtcList")
	commonRep:=crv.CommonReq{
		ModelID:"diag_manual_fault",
		Filter:&map[string]interface{}{
			"ecu_id":&map[string]interface{}{
				"Op.in":dtc.Ecus,
			},
		},
		Fields:&QueryDTCSignalFields,
	}

	return dtc.CRVClient.Query(&commonRep)
}

func (dtc *dtcList)convertDtcList(queryResult *crv.CommonRsp)([]interface{},int){
	log.Println("start convertDtcList")
	list,ok:=queryResult.Result["list"]
	if !ok {
		log.Println("convertDtcList queryResult no list")
		return nil,common.ResultNoDtcList
	}

	dtcList,ok:=list.([]interface{})
	if !ok || len(dtcList)<=0 {
		log.Println("convertDtcList queryResult no list")
		return nil,common.ResultNoDtcList
	}

	for _,item:=range dtcList {
		mapItem:=item.(map[string]interface{})
		mapItem["DtcNum"]=mapItem["dtc_id"]
		signals,_:=mapItem["signals"].(map[string]interface{})["list"].([]interface{})
		for _,signal:=range signals {
			signalItem:=signal.(map[string]interface{})
			sCanID,_:=signalItem["can_id"].(string)
			sPduID,_:=signalItem["pdu_id"].(string)
			sStartAddr,_:=signalItem["start_addr"].(string)
			signalItem["SignalID"]=fmt.Sprintf("%s:%s:%s",sCanID,sPduID,sStartAddr)
			//signalItem["SignalID"]=signalItem["id"]
			signalItem["SignalName"]=signalItem["name"]
			delete(signalItem,"id")
			delete(signalItem,"name")
			delete(signalItem,"can_id")
			delete(signalItem,"pdu_id")
			delete(signalItem,"start_addr")
		}
		mapItem["CorrelationSignal"]=signals
		delete(mapItem,"id")
		delete(mapItem,"dtc_id")
		delete(mapItem,"signals")
	}	

	return dtcList,common.ResultSuccess
}