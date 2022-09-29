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
	{"field": "ecu_id"},
	{
		"field": "signals",
		"fieldType": "many2many",
		"relatedModelID": "diag_signal",
		"pagination":map[string]interface{}{
			"pageSize":10000,
			"current":1,
		},
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		 	{"field": "name"},
			{"field": "can_id"},
			{"field": "start_addr"},
			{"field": "pdu_id"},
			{"field":"byte_order"},
			{"field":"len"},
		},
	},
}

type dtcList struct {
	CRVClient *crv.CRVClient
	Ecus []string
	EcuChannelMap map[string]interface{}
	DtcList []interface{}
}

func getDtcList(
	ecus []string,
	channelMap map[string]interface{},
	crvClient *crv.CRVClient,
	signalList *map[string]interface{})(*dtcList,int){
	dtc:=&dtcList{
		CRVClient:crvClient,
		Ecus:ecus,
		EcuChannelMap:channelMap,
	}

	log.Println("getDtcList",ecus)

	rsp,errorCode:=dtc.queryDtcList()
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	dtc.DtcList,errorCode=dtc.convertDtcList(rsp,signalList)
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
		Pagination:&crv.Pagination{
			PageSize:5000,
			Current:1,
		},
		Fields:&QueryDTCSignalFields,
	}

	return dtc.CRVClient.Query(&commonRep)
}

func (dtc *dtcList)convertDtcList(
	queryResult *crv.CommonRsp,
	signalList *map[string]interface{})([]interface{},int){
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
		ecuID,_:=mapItem["ecu_id"].(string)
		channel,_:=dtc.EcuChannelMap[ecuID]
		signals,_:=mapItem["signals"].(map[string]interface{})["list"].([]interface{})
		for _,signal:=range signals {
			signalItem:=signal.(map[string]interface{})
			sCanID,_:=signalItem["can_id"].(string)
			sPduID,_:=signalItem["pdu_id"].(string)
			sType,_:=signalItem["byte_order"].(string)
			sStartAddr,_:=signalItem["start_addr"].(string)
			signalID:=fmt.Sprintf("%s:%s:%s:%s",sCanID,sPduID,sStartAddr,channel)
			signalItem["SignalID"]=	signalID		
		
			//判断并添加signal
			_,ok:=(*signalList)[signalID]
			if !ok {
				(*signalList)[signalID]=map[string]interface{}{
					"Channel":channel,
					"CanID":sCanID,
					"Type":sType,
					"SignalName":signalItem["name"],
					"PduId":sPduID,
					"startAddr":sStartAddr,
					"len":signalItem["len"],
				}
			}
			
			//signalItem["SignalID"]=signalItem["id"]
			signalItem["SignalName"]=signalItem["name"]
			delete(signalItem,"id")
			delete(signalItem,"name")
			delete(signalItem,"can_id")
			delete(signalItem,"pdu_id")
			delete(signalItem,"start_addr")
			delete(signalItem,"len")
			delete(signalItem,"byte_order")
		}
		mapItem["CorrelationSignal"]=signals
		delete(mapItem,"id")
		delete(mapItem,"dtc_id")
		delete(mapItem,"signals")
		delete(mapItem,"ecu_id")
	}	

	return dtcList,common.ResultSuccess
}