package send

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"	
)

var QuerySendEventRecordFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "device_number"},
	{"field": "version"},
}

func getSaveEventRecord(
	vehicle sendVehicleItem,
	sendUser string,
	parameter string)(map[string]interface{}){
	row:=map[string]interface{}{
		"device_number":vehicle.DeviceID,
		"vin":vehicle.Vin,
		"status":"0",
		"test_specification":vehicle.TestSpecification,
		"platform_id":vehicle.ProjectNum,
		"paramter":parameter,
		"send_user":sendUser,
		crv.SAVE_TYPE_COLUMN:crv.SAVE_CREATE,
	}
	return row
}

func saveSendEventRecords(
	crvClient *crv.CRVClient,
 	rows []map[string]interface{},
	 token string)(*crv.CommonRsp,int){
	
	commonRep:=crv.CommonReq{
		ModelID:"diag_event_sendrecord",
		List:&rows,
	}

	return crvClient.Save(&commonRep,token)
}

func getSendEventRecords(
	crvClient *crv.CRVClient,
	sendRsp *crv.CommonRsp,
	token string)(*crv.CommonRsp,int){

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
		ModelID:"diag_event_sendrecord",
		Filter:&map[string]interface{}{
			"id":&map[string]interface{}{
				"Op.in":ids,
			},
		},
		Fields:&QuerySendEventRecordFields,
	}

	return crvClient.Query(&commonRep,token)
}

func cacheSendEventRecords(
	sendRecordCache *SendRecordCache,
	sendRsp *crv.CommonRsp)(int){
	list,_:=sendRsp.Result["list"]
	records,_:=list.([]interface{})
	log.Println("cacheSendEventRecords")
	log.Println(records)
	for _,row:=range(records) {
		rowMap,_:=row.(map[string]interface{})
		device,_:=rowMap["device_number"]
		deviceID,_:=device.(string)
		log.Println(deviceID)
		err:=sendRecordCache.SaveSendRecord("event_"+deviceID,rowMap)
		if err!=nil {
			return common.ResultCacheSendRecError
		}
	}

	return common.ResultSuccess
}

func createSendEventRecords(
	crvClient *crv.CRVClient,
	sendRecordCache *SendRecordCache,
	vehicles []sendVehicleItem,
	sendUser string,
	parameter string,
	token string)(*crv.CommonRsp,int){
	
	saveList:=make([]map[string]interface{},len(vehicles))
	for index,vehicle:=range(vehicles){
		saveList[index]=getSaveEventRecord(vehicle,sendUser,parameter)
	}

	rsp,errorCode:=saveSendEventRecords(crvClient,saveList,token)
	if errorCode!=common.ResultSuccess {
		return rsp,errorCode
	}

	rsp,errorCode=getSendEventRecords(crvClient,rsp,token)
	if errorCode!=common.ResultSuccess {
		return rsp,errorCode
	}

	errorCode=cacheSendEventRecords(sendRecordCache,rsp)

	return nil,errorCode
}