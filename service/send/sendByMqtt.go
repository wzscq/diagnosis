package send

import (
	"digimatrix.com/diagnosis/mqtt"
	"digimatrix.com/diagnosis/common"
	"log"	
	"encoding/json"
)

func sendByMqtt(
	mqttClient *mqtt.MQTTClient,
	vehicles []sendVehicleItem,
	parameter string)(int){
	
	log.Println("start sendByMqtt")
	for _,vehicle:=range(vehicles){
		topic:="MQB/"+vehicle.DeviceID+"/Diag"
		log.Println(topic)
		errorCode:=mqttClient.Publish(topic,parameter)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
	}

	log.Println("end sendByMqtt")
	return common.ResultSuccess
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

func getVehicleDiagSignals(
	deviceID string,
	signalList *map[string]interface{},
	deviceSignalCache *DeviceSignalCache)(string){
	eventSigLst,_:=deviceSignalCache.GetSignalList(deviceID,"event")
	sendSigList:=signalList
	if eventSigLst!=nil {
		sendSigList=mergeSignalList(signalList,eventSigLst)
	}
	deviceSignalCache.SaveSignalList(deviceID,"diag",signalList)
	return convertToSignalParameter(sendSigList)
}

func sendDiagSignalByMqtt(
	mqttClient *mqtt.MQTTClient,
	vehicles []sendVehicleItem,
	signalList *map[string]interface{},
	deviceSignalCache *DeviceSignalCache)(int){
	
	log.Println("start sendDiagSignalByMqtt")
	for _,vehicle:=range(vehicles){
		parameter:=getVehicleDiagSignals(vehicle.DeviceID,signalList,deviceSignalCache)
		topic:="MQB/"+vehicle.DeviceID+"/SignalFilter"
		log.Println(topic)
		log.Println(parameter)
		errorCode:=mqttClient.Publish(topic,parameter)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
	}

	log.Println("end sendDiagSignalByMqtt")
	return common.ResultSuccess
} 