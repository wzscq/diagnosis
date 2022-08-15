package send

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"	
)

var QuerySendRecordFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "device_number"},
	{"field": "version"},
}

func getSaveRecord(
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

func saveSendRecords(
	crvClient *crv.CRVClient,
 	rows []map[string]interface{})(*crv.CommonRsp,int){
	
	commonRep:=crv.CommonReq{
		ModelID:"diag_param_sendrecord",
		List:&rows,
	}

	return crvClient.Save(&commonRep)
}

func getSendRecords(
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
		ModelID:"diag_param_sendrecord",
		Filter:&map[string]interface{}{
			"id":&map[string]interface{}{
				"Op.in":ids,
			},
		},
		Fields:&QuerySendRecordFields,
	}

	return crvClient.Query(&commonRep)
}

func cacheSendRecords(
	sendRecordCache *SendRecordCache,
	sendRsp *crv.CommonRsp)(int){
	list,_:=sendRsp.Result["list"]
	records,_:=list.([]interface{})
	log.Println("cacheSendRecords")
	log.Println(records)
	for _,row:=range(records) {
		rowMap,_:=row.(map[string]interface{})
		device,_:=rowMap["device_number"]
		deviceID,_:=device.(string)
		log.Println(deviceID)
		err:=sendRecordCache.SaveSendRecord(deviceID,rowMap)
		if err!=nil {
			return common.ResultCacheSendRecError
		}
	}

	return common.ResultSuccess
}

func createSendRecords(
	crvClient *crv.CRVClient,
	sendRecordCache *SendRecordCache,
	vehicles []sendVehicleItem,
	sendUser string,
	parameter string)(*crv.CommonRsp,int){
	
	saveList:=make([]map[string]interface{},len(vehicles))
	for index,vehicle:=range(vehicles){
		saveList[index]=getSaveRecord(vehicle,sendUser,parameter)
	}

	rsp,errorCode:=saveSendRecords(crvClient,saveList)
	if errorCode!=common.ResultSuccess {
		return rsp,errorCode
	}

	rsp,errorCode=getSendRecords(crvClient,rsp)
	if errorCode!=common.ResultSuccess {
		return rsp,errorCode
	}

	errorCode=cacheSendRecords(sendRecordCache,rsp)

	return nil,errorCode
}