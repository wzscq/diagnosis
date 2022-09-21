package send

import (
	"digimatrix.com/diagnosis/mqtt"
	"digimatrix.com/diagnosis/common"
	"digimatrix.com/diagnosis/crv"
	"log"	
	"encoding/json"
)

var QuerySendSignalRecordFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "device_number"},
	{"field": "version"},
}

func getSaveSignalRecord(
	vehicle sendVehicleItem,
	sendUser string,
	parameter string)(map[string]interface{}){
	row:=map[string]interface{}{
		"device_number":vehicle.DeviceID,
		"vin":vehicle.Vin,
		"status":"0",
		"paramter":parameter,
		"send_user":sendUser,
		crv.SAVE_TYPE_COLUMN:crv.SAVE_CREATE,
	}
	return row
}

func saveSendSignalRecord(
	crvClient *crv.CRVClient,
 	rows []map[string]interface{})(*crv.CommonRsp,int){
	
	commonRep:=crv.CommonReq{
		ModelID:"diag_signal_sendrecord",
		List:&rows,
	}

	return crvClient.Save(&commonRep)
}

func getSendSignalRecord(
	crvClient *crv.CRVClient,
	sendRsp *crv.CommonRsp)(*crv.CommonRsp,int){

	list,_:=sendRsp.Result["list"]
	records,_:=list.([]interface{})
	
	ids:=[]interface{}{}
	for _,row:=range(records) {
		rowMap,_:=row.(map[string]interface{})
		log.Println(rowMap)
		id,_:=rowMap["id"]	
		ids=append(ids,id)
	}
	log.Println(ids)
	commonRep:=crv.CommonReq{
		ModelID:"diag_signal_sendrecord",
		Filter:&map[string]interface{}{
			"id":&map[string]interface{}{
				"Op.in":ids,
			},
		},
		Fields:&QuerySendSignalRecordFields,
	}

	return crvClient.Query(&commonRep)
}

func cacheSendSignalRecord(
	sendRecordCache *SendRecordCache,
	sendRsp *crv.CommonRsp)(int){
	list,_:=sendRsp.Result["list"]
	records,_:=list.([]interface{})
	log.Println("cacheSendSignalRecord")
	log.Println(records)
	for _,row:=range(records) {
		rowMap,_:=row.(map[string]interface{})
		device,_:=rowMap["device_number"]
		deviceID,_:=device.(string)
		log.Println(deviceID)
		err:=sendRecordCache.SaveSendRecord("signal_"+deviceID,rowMap)
		if err!=nil {
			return common.ResultCacheSendRecError
		}
	}

	return common.ResultSuccess
}

func createSendSignalRecord(
	crvClient *crv.CRVClient,
	sendRecordCache *SendRecordCache,
	vehicle sendVehicleItem,
	sendUser string,
	parameter string)(*crv.CommonRsp,int){
	
	saveList:=make([]map[string]interface{},1)
	saveList[0]=getSaveSignalRecord(vehicle,sendUser,parameter)
	
	rsp,errorCode:=saveSendSignalRecord(crvClient,saveList)
	if errorCode!=common.ResultSuccess {
		log.Println("createSendSignalRecord saveSendSignalRecord error")
		return rsp,errorCode
	}

	rsp,errorCode=getSendSignalRecord(crvClient,rsp)
	if errorCode!=common.ResultSuccess {
		log.Println("createSendSignalRecord getSendSignalRecord error")
		return rsp,errorCode
	}

	errorCode=cacheSendSignalRecord(sendRecordCache,rsp)

	return nil,errorCode
}

func mergeSignalList(
	signalDiag,signalEvent *map[string]interface{})(*map[string]interface{}){
	result:=map[string]interface{}{}
	
	for key,sigItem:=range(*signalDiag){
		result[key]=sigItem
	}

	for key,sigItem:=range(*signalEvent){
		_,ok:=result[key]
		if !ok {
			result[key]=sigItem
		}
	}
	return &result
}

func convertToSignalParameter(signalList *map[string]interface{})(string){
	canIDList:=[]map[string]interface{}{}
	for _,sigItem:=range(*signalList) {
		sigItemMap,_:=sigItem.(map[string]interface{})
		findID:=-1
		for index,canItem:=range(canIDList){
			if(sigItemMap["CanID"]==canItem["CanID"] && sigItemMap["Channel"]==canItem["Channel"]){
				findID=index
				break
			}
		}
		canSigItem:=map[string]interface{}{
			"SignalName":sigItemMap["SignalName"],
			"PduId":sigItemMap["PduId"],
			"startAddr":sigItemMap["startAddr"],
			"len":sigItemMap["len"],
		}
		if findID<0 {
			canItem:=map[string]interface{}{
				"CanID":sigItemMap["CanID"],
				"Channel":sigItemMap["Channel"],
				"Type":sigItemMap["Type"],
				"SignalNameList":[]interface{}{
					canSigItem,	
				},
			}
			canIDList=append(canIDList,canItem)
		} else {
			canSigLst:=canIDList[findID]["SignalNameList"].([]interface{})
			canSigLst=append(canSigLst,canSigItem)
			canIDList[findID]["SignalNameList"]=canSigLst
		}
	}	
	mapSigParam:=map[string]interface{}{
		"CanIDList":canIDList,
	}
	bytes, err := json.Marshal(mapSigParam)
	if err!=nil {
		log.Println("convertToSignalParameter error:",err.Error())
		return ""
	}
    // Convert bytes to string.
    jsonStr := string(bytes)
	return jsonStr
}

func getVehicleSignals(
	deviceID string,
	signalList *map[string]interface{},
	deviceSignalCache *DeviceSignalCache,
	busi string)(string){
	busiOther:="diag"
	if busi=="diag"{
		busiOther="event"
	}
	log.Printf("getVehicleSignals deviceID:%s,busi:%s",deviceID,busiOther)
	otherSigLst,_:=deviceSignalCache.GetSignalList(deviceID,busiOther)
	sendSigList:=signalList
	if otherSigLst!=nil {
		sendSigList=mergeSignalList(signalList,otherSigLst)
	}
	
	err:=deviceSignalCache.SaveSignalList(deviceID,busi,signalList)
	if err!=nil {
		log.Println(err)
	}
	return convertToSignalParameter(sendSigList)
}

func sendSignalByMqtt(
	mqttClient *mqtt.MQTTClient,
	vehicles []sendVehicleItem,
	signalList *map[string]interface{},
	deviceSignalCache *DeviceSignalCache,
	busi string,
	crvClient *crv.CRVClient,
	sendRecordCache *SendRecordCache,
	sendUser string)(int){
	
	log.Println("start sendSignalByMqtt")
	for _,vehicle:=range(vehicles){
		parameter:=getVehicleSignals(vehicle.DeviceID,signalList,deviceSignalCache,busi)
		topic:="MQB/"+vehicle.DeviceID+"/SignalFilter"
		log.Println(topic)
		log.Println(parameter)
		errorCode:=mqttClient.Publish(topic,parameter)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
		createSendSignalRecord(
			crvClient,
			sendRecordCache,
			vehicle,
			sendUser,
			parameter)
	}

	log.Println("end sendSignalByMqtt")
	return common.ResultSuccess
} 