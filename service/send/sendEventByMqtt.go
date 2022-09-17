package send

import (
	"digimatrix.com/diagnosis/mqtt"
	"digimatrix.com/diagnosis/common"
	"log"	
)

func sendEventByMqtt(
	mqttClient *mqtt.MQTTClient,
	vehicles []sendVehicleItem,
	parameter string)(int){
	
	log.Println("start sendEventByMqtt")
	for _,vehicle:=range(vehicles){
		topic:="MQB/"+vehicle.DeviceID+"/Event"
		log.Println(topic)
		errorCode:=mqttClient.Publish(topic,parameter)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
	}

	log.Println("end sendByMqtt")
	return common.ResultSuccess
}

func getVehicleEventSignals(
	deviceID string,
	signalList *map[string]interface{},
	deviceSignalCache *DeviceSignalCache)(string){
	diagSigLst,_:=deviceSignalCache.GetSignalList(deviceID,"diag")
	sendSigList:=signalList
	if diagSigLst!=nil {
		sendSigList=mergeSignalList(signalList,diagSigLst)
	}
	deviceSignalCache.SaveSignalList(deviceID,"event",signalList)
	return convertToSignalParameter(sendSigList)
}

func sendEventSignalByMqtt(
	mqttClient *mqtt.MQTTClient,
	vehicles []sendVehicleItem,
	signalList *map[string]interface{},
	deviceSignalCache *DeviceSignalCache)(int){
	
	log.Println("start sendEventSignalByMqtt")
	for _,vehicle:=range(vehicles){
		parameter:=getVehicleEventSignals(vehicle.DeviceID,signalList,deviceSignalCache)
		topic:="MQB/"+vehicle.DeviceID+"/SignalFilter"
		log.Println(topic)
		log.Println(parameter)
		errorCode:=mqttClient.Publish(topic,parameter)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
	}

	log.Println("end sendEventSignalByMqtt")
	return common.ResultSuccess
}