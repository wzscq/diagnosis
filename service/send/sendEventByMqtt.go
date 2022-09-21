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
