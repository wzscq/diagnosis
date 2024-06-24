package send

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
	"fmt"
	"strconv"
	"strings"
)

var QueryParameterFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "version"},
	{"field": "name"},
	{"field": "platform_id"},
	{"field": "time_offset"},
	{"field": "channel"},
	{"field": "use_triggercanid"},
	/*{
		"field": "domain_id",
		"fieldType": "MANY_TO_ONE",
		"relatedModelID": "diag_domain",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "name"},
		},
	},*/
	{
		"field": "ecu_id",
		"fieldType": "many2one",
		"relatedModelID": "diag_ecu",
		"pagination":map[string]interface{}{
			"pageSize":5000,
			"current":1,
		},
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		 	{"field": "name"},
			{"field": "tx"},
			{"field": "rx"},
			{"field": "instruct"},
			{"field": "trigger_can_id"},
		},
	},
	{
		"field": "logistics_did",
		"fieldType": "many2many",
		"relatedModelID": "diag_logistics",
		"associationModelID": "diag_parameter_logistics_did",
		"pagination":map[string]interface{}{
			"pageSize":5000,
			"current":1,
		},
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "did"},
		},
	},
	{
		"field": "internal_did",
		"fieldType": "many2many",
		"relatedModelID": "diag_logistics",
		"associationModelID": "diag_parameter_internal_did",
		"pagination":map[string]interface{}{
			"pageSize":5000,
			"current":1,
		},
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		  	{"field": "did"},
		},
	},
	{
		"field": "triggers",
		"fieldType": "one2many",
		"relatedModelID": "diag_parameter_trigger",
		"relatedField": "parameter_id",
		"pagination":map[string]interface{}{
			"pageSize":5000,
			"current":1,
		},
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		  	{
				"field": "diag_signal_id",
				"fieldType": "many2one",
				"relatedModelID": "diag_signal",
				"pagination":map[string]interface{}{
					"pageSize":5000,
					"current":1,
				},
				"fields": []map[string]interface{}{
			  		{"field": "id"},
			  		{"field": "name"},
					{"field": "can_id"},
					{"field": "start_addr"},
					{"field": "pdu_id"},
					{"field": "factor"},
					{"field": "offset"},
					{"field":"byte_order"},
					{"field":"len"},
				},
		  	},
		  	{"field": "logic"},
		  	{"field": "value"},
		  	//{"field": "version"},
		  	{"field": "parameter_id"},
		},
	},
	//{"field": "remark"},
}

type diagParameter struct {
	CRVClient *crv.CRVClient
	Ids []string
	Records []map[string]interface{}
	Params []map[string]interface{}
	PlatformID string
}

func getDiagParams(
	ids []string,
	crvClient *crv.CRVClient,
	signalList *map[string]interface{},
	token string)(*diagParameter,int){
	dp:=&diagParameter{
		CRVClient:crvClient,
		Ids:ids,
	}

	//获取诊断参数配置数据
	rsp,errorCode:=dp.queryDiagParameter(token)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}
	//将结果转换为map数组，存放在Records成员中
	errorCode=dp.updateRecords(rsp.Result)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}
	//检查选择参数是否符合规则
	errorCode=dp.isValid()
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	dp.convertDiagParameters(signalList)
	return dp,common.ResultSuccess
}

func (dp *diagParameter)updateRecords(result map[string]interface{})(int){
	list,ok:=result["list"]
	if !ok {
		log.Println("getDiagParams queryResult no list")
		return common.ResultNoParams
	}

	records,ok:=list.([]interface{})
	if !ok || len(records)<=0 {
		log.Println("getDiagParams queryResult no list")
		return common.ResultNoParams
	}

	dp.Records=make([]map[string]interface{},len(records))
	for index,row:=range(records){
		dp.Records[index]=row.(map[string]interface{})
	}
	return common.ResultSuccess
}

func (dp *diagParameter)queryDiagParameter(token string)(*crv.CommonRsp,int){
	log.Println("start queryDiagParameter")
	commonRep:=crv.CommonReq{
		ModelID:"diag_parameter",
		Filter:&map[string]interface{}{
			"id":&map[string]interface{}{
				"Op.in":dp.Ids,
			},
		},
		Fields:&QueryParameterFields,
	}

	return dp.CRVClient.Query(&commonRep,token)
}

func (dp *diagParameter)isValid()(int){
	//下发参数的项目号要求一致，ECU不能重复
	mapProject:=map[string]interface{}{}
	mapEcu:=map[string]interface{}{}
	for _,record:=range(dp.Records){
		prj,_:=record["platform_id"]
		prjStr,_:=prj.(string)
		dp.PlatformID=prjStr
		mapProject[prjStr]=prj
		
		ecuList,ok:=record["ecu_id"].(map[string]interface{})["list"].([]interface{})
		if !ok || len(ecuList)==0 {
			return common.ResultParamWithoutEcu
		}

		ecu,_:=ecuList[0].(map[string]interface{})["id"]
		ecuStr,_:=ecu.(string)
		_,ok=mapEcu[ecuStr]
		if ok {
			return common.ResultRepeatedEcu
		}
		mapEcu[ecuStr]=ecu
	}

	if len(mapProject)!=1 {
		return common.ResultMultiProject
	}

	return common.ResultSuccess
}

func (dp *diagParameter)convertDiagParameter(
	row map[string]interface{},
	sendSignalList *map[string]interface{})(map[string]interface{}){
	diagParaList:=map[string]interface{}{}
	diagParaList["TimeOffset"]=row["time_offset"]
	diagParaList["Channel"]=row["channel"]

	use_triggercanid,_:=row["use_triggercanid"].(string)
	
	ecuList,ok:=row["ecu_id"].(map[string]interface{})["list"].([]interface{})
	if ok && len(ecuList)>0 {
		ecu,_:=ecuList[0].(map[string]interface{})
		diagParaList["Ecu"]=ecu["id"]
		diagParaList["TxId"]=ecu["tx"]
		diagParaList["RxId"]=ecu["rx"]
		if use_triggercanid=="1" {
			diagParaList["TriggerCanId"]=ecu["trigger_can_id"]
		} else {
			diagParaList["TriggerCanId"]="0"
		}
		diagParaList["DiagInstruct"]=ecu["instruct"]
	}
	
	logistics,_:=row["logistics_did"].(map[string]interface{})["list"].([]interface{})
	for _,item:=range logistics {
		mapItem:=item.(map[string]interface{})
		delete(mapItem,"id")
		mapItem["DID"]=mapItem["did"]
		delete(mapItem,"did")
	}

	internallogistics,_:=row["internal_did"].(map[string]interface{})["list"].([]interface{})
	for _,item:=range internallogistics {
		mapItem:=item.(map[string]interface{})
		delete(mapItem,"id")
		mapItem["DID"]=mapItem["did"]
		delete(mapItem,"did")
	}

	var triggerList []interface{}
	if val, ok := row["triggers"]; ok {
		triggerList=val.(map[string]interface{})["list"].([]interface{})
	} else {
		triggerList=[]interface{}{}
	}
	
	for _,item:=range triggerList {
		mapItem:=item.(map[string]interface{})
		delete(mapItem,"id")
		delete(mapItem,"parameter_id")
		signalList,ok:=mapItem["diag_signal_id"].(map[string]interface{})["list"].([]interface{})
		var floatFactor,floatOffset float64
		var err error
		if ok && len(signalList)>0 {
			signal,_:=signalList[0].(map[string]interface{})
			mapItem["SignalName"]=signal["name"]
			
			sCanID,_:=signal["can_id"].(string)
			sPduID,_:=signal["pdu_id"].(string)
			sType,_:=signal["byte_order"].(string)
			sStartAddr,_:=signal["start_addr"].(string)
			signalID:=fmt.Sprintf("%s:%s:%s:%s",sCanID,sPduID,sStartAddr,diagParaList["Channel"])
			mapItem["SignalID"]=signalID
			//判断并添加signal
			_,ok:=(*sendSignalList)[signalID]
			if !ok {
				(*sendSignalList)[signalID]=map[string]interface{}{
					"Channel":diagParaList["Channel"],
					"CanID":sCanID,
					"Type":sType,
					"SignalName":signal["name"],
					"PduId":sPduID,
					"startAddr":sStartAddr,
					"len":signal["len"],
				}
			}
			sFactor,_:=signal["factor"].(string)
			floatFactor, err = strconv.ParseFloat(sFactor, 64)
			if err!=nil {
				log.Println("trigger signal factor can not convert to float64")
				log.Println("factor value is :"+sFactor)
				return nil
			}
			sOffset,_:=signal["offset"].(string)
			floatOffset, err= strconv.ParseFloat(sOffset, 64)
			if err!=nil {
				log.Println("trigger signal offset can not convert to float64")
				return nil
			}
		}
		delete(mapItem,"diag_signal_id")
		mapItem["Logic"]=mapItem["logic"]
		delete(mapItem,"logic")
		sValue,_:=mapItem["value"].(string)
		floatVal, err := strconv.ParseFloat(sValue, 64)
		if err!=nil {
			log.Println("trigger value can not convert to float64")
			return nil
		}
		//mapItem["Value"]=fmt.Sprintf("%.f",floatVal/floatFactor+floatOffset)
		mapItem["Value"]=fmt.Sprintf("%.f",(floatVal-floatOffset)/floatFactor)
		delete(mapItem,"value")
	}
	
	diagParaList["Logistics"]=logistics
	diagParaList["Internallogistics"]=internallogistics
	diagParaList["TriggerList"]=triggerList

	return diagParaList
}

func (dp *diagParameter)getLogic(logic interface{})(string){
	strLogic,_:=logic.(string)
	log.Println("diagParameter getLogic strLogic before replace:"+strLogic)
	//替换<,>符号
	strLogic=strings.Replace(strLogic,"\u003e",">",1)
	strLogic=strings.Replace(strLogic,"\u003c","<",1)
	log.Println("diagParameter getLogic strLogic after replace:"+strLogic)
	return strLogic
}

func (dp *diagParameter)convertDiagParameters(signalList *map[string]interface{}){
	log.Println("start convertDiagParameters")
	dp.Params=make([]map[string]interface{},len(dp.Records))
	
	for index,row:=range(dp.Records){
		dp.Params[index]=dp.convertDiagParameter(row,signalList)
	}

	log.Println("end convertDiagParameters")
} 

func (dp *diagParameter)getProjectNum()(interface{}){
	return dp.Records[0]["platform_id"]
}

func (dp *diagParameter)getEcuIDs()([]string){
	ids:=[]string{}
	for _,param:=range(dp.Params){
		id,_:=param["Ecu"]
		strID,_:=id.(string)
		ids=append(ids,strID)
	}
	return ids
}

func (dp *diagParameter)getEcuChannelMap()(map[string]interface{}){
	ecuChanelMap:=map[string]interface{}{}
	for _,param:=range(dp.Params){
		id,_:=param["Ecu"]
		strID,_:=id.(string)
		ecuChanelMap[strID]=param["Channel"]
	}
	return ecuChanelMap
}
