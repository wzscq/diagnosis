package send

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
	"fmt"
	"strconv"
	"strings"
)

var QueryEventParameterFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "version"},
	{"field": "name"},
	{"field": "platform_id"},
	{"field": "last_time"},
	{"field": "next_time"},
	{"field": "domain_id"},
	{"field": "channel"},
	{
		"field": "triggers",
		"fieldType": "one2many",
		"relatedModelID": "diag_event_parameter_trigger",
		"relatedField": "event_parameter_id",
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		  	{
				"field": "diag_signal_id",
				"fieldType": "many2one",
				"relatedModelID": "diag_signal",
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
		  	{"field": "event_parameter_id"},
		},
	},
	{
		"field": "correlation_signals",
		"fieldType": "many2many",
		"relatedModelID": "diag_signal",
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
	//{"field": "remark"},
}

type eventParameter struct {
	CRVClient *crv.CRVClient
	Ids []string
	Records []map[string]interface{}
	Params []map[string]interface{}
}

func getEventParams(
	ids []string,
	crvClient *crv.CRVClient,
	sendSignalList *map[string]interface{})(*eventParameter,int){
	dp:=&eventParameter{
		CRVClient:crvClient,
		Ids:ids,
	}

	//获取诊断参数配置数据
	rsp,errorCode:=dp.queryEventParameter()
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

	dp.convertDiagParameters(sendSignalList)
	return dp,common.ResultSuccess
}

func (dp *eventParameter)updateRecords(result map[string]interface{})(int){
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

func (dp *eventParameter)queryEventParameter()(*crv.CommonRsp,int){
	log.Println("start queryEventParameter")
	commonRep:=crv.CommonReq{
		ModelID:"diag_event_parameter",
		Filter:&map[string]interface{}{
			"id":&map[string]interface{}{
				"Op.in":dp.Ids,
			},
		},
		Fields:&QueryEventParameterFields,
	}

	return dp.CRVClient.Query(&commonRep)
}

func (dp *eventParameter)isValid()(int){
	//下发参数的项目号要求一致，ECU不能重复
	mapProject:=map[string]interface{}{}
	for _,record:=range(dp.Records){
		prj,_:=record["platform_id"]
		prjStr,_:=prj.(string)
		mapProject[prjStr]=prj
	}

	if len(mapProject)!=1 {
		return common.ResultMultiProject
	}

	return common.ResultSuccess
}

func (dp *eventParameter)convertDiagParameter(
	row map[string]interface{},
	sendSignalList *map[string]interface{})(map[string]interface{}){
	diagParaList:=map[string]interface{}{}
	diagParaList["EventID"]=row["id"]
	diagParaList["EventName"]=row["name"]
	diagParaList["TriggerLogic"]=""
	diagParaList["LastTime"]=row["last_time"]
	diagParaList["NextTime"]=row["next_time"]
	
	triggerMap,ok:=row["triggers"].(map[string]interface{})
	if ok {
		triggerList,_:=triggerMap["list"].([]interface{})
		for _,item:=range triggerList {
			mapItem:=item.(map[string]interface{})
			delete(mapItem,"id")
			delete(mapItem,"event_parameter_id")
			signalList,ok:=mapItem["diag_signal_id"].(map[string]interface{})["list"].([]interface{})
			var floatFactor,floatOffset float64
			var err error
			if ok && len(signalList)>0 {
				signal,_:=signalList[0].(map[string]interface{})
				mapItem["SignalName"]=signal["name"]
				
				sCanID,_:=signal["can_id"].(string)
				sPduID,_:=signal["pdu_id"].(string)
				sStartAddr,_:=signal["start_addr"].(string)

				signalID:=fmt.Sprintf("%s:%s:%s:%s",sCanID,sPduID,sStartAddr,row["channel"])
				mapItem["SignalID"]=signalID
				//判断并添加signal
				_,ok:=(*sendSignalList)[signalID]
				if !ok {
					(*sendSignalList)[signalID]=map[string]interface{}{
						"Channel":row["channel"],
						"CanID":sCanID,
						"Type":signal["byte_order"],
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
			mapItem["Value2"]=fmt.Sprintf("%.f",floatVal/floatFactor+floatOffset)
			mapItem["Value1"]=""
			delete(mapItem,"value")
		}
		diagParaList["Trigger"]=triggerList
	} else {
		diagParaList["Trigger"]=[]interface{}{}
	}

	signalMap,ok:=row["correlation_signals"].(map[string]interface{})
	if ok {
		signals,_:=signalMap["list"].([]interface{})
		for _,signal:=range signals {
			signalItem:=signal.(map[string]interface{})
			sCanID,_:=signalItem["can_id"].(string)
			sPduID,_:=signalItem["pdu_id"].(string)
			sStartAddr,_:=signalItem["start_addr"].(string)
			
			
			signalID:=fmt.Sprintf("%s:%s:%s:%s",sCanID,sPduID,sStartAddr,row["channel"])
			signalItem["SignalID"]=signalID
			//判断并添加signal
			_,ok:=(*sendSignalList)[signalID]
			if !ok {
				(*sendSignalList)[signalID]=map[string]interface{}{
					"Channel":row["channel"],
					"CanID":sCanID,
					"Type":signalItem["byte_order"],
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
			delete(signalItem,"byte_order")
			delete(signalItem,"len")
		}
		diagParaList["CorrelationSignal"]=signals
	} else {
		diagParaList["CorrelationSignal"]=[]interface{}{}
	}
	
	return diagParaList
}

func (dp *eventParameter)getLogic(logic interface{})(string){
	strLogic,_:=logic.(string)
	log.Println("diagParameter getLogic strLogic before replace:"+strLogic)
	//替换<,>符号
	strLogic=strings.Replace(strLogic,"\u003e",">",1)
	strLogic=strings.Replace(strLogic,"\u003c","<",1)
	log.Println("diagParameter getLogic strLogic after replace:"+strLogic)
	return strLogic
}

func (dp *eventParameter)convertDiagParameters(sendSignalList *map[string]interface{}){
	log.Println("start convertDiagParameters")
	dp.Params=make([]map[string]interface{},len(dp.Records))
	
	for index,row:=range(dp.Records){
		dp.Params[index]=dp.convertDiagParameter(row,sendSignalList)
	}

	log.Println("end convertDiagParameters")
} 

func (dp *eventParameter)getProjectNum()(interface{}){
	return dp.Records[0]["platform_id"]
}

func (dp *eventParameter)getEcuIDs()([]string){
	ids:=[]string{}
	for _,param:=range(dp.Params){
		id,_:=param["Ecu"]
		strID,_:=id.(string)
		ids=append(ids,strID)
	}
	return ids
}
