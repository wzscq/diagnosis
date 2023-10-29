package busi

import (
	"log"
	"digimatrix.com/diagnosis/crv"
)

const (
	MODEL_HEARTBEAT = "diag_device_heartbeat"
	MODEL_VEHICLE  = "vehiclemanagement"
)

func getDeviceHeartbeatSaveReq(deviceID,vin string)(*crv.CommonReq){
	return &crv.CommonReq{
		ModelID:MODEL_HEARTBEAT,
		List:&[]map[string]interface{}{
			map[string]interface{}{
				"device_number":deviceID,
				"vin":vin,
				crv.SAVE_TYPE_COLUMN:crv.SAVE_CREATE,
			},
		},
	}
}

func getUpdateVinDeviceReq(deviceID,vin string,version interface{})(*crv.CommonReq){
	return &crv.CommonReq{
		ModelID:MODEL_VEHICLE,
		List:&[]map[string]interface{}{
			map[string]interface{}{
				"id":vin,
				"DeviceNumber":deviceID,
				crv.CC_VERSION:version,
				crv.SAVE_TYPE_COLUMN:crv.SAVE_UPDATE,
			},
		},
	}
}

func getQueryVinReq(vin string)(*crv.CommonReq){
	return &crv.CommonReq{
		ModelID:MODEL_VEHICLE,
		Fields:&[]map[string]interface{}{
			{"field": "id"},
			{"field": "version"},
		},
		Filter:&map[string]interface{}{
			"id":vin,
		},
	}
}

func getVersion(result map[string]interface{})(interface{}){
	list,ok:=result["list"]
	if !ok {
		log.Println("busiHeartbeat getVersion no list")
		return nil
	}

	listMap,ok:=list.([]interface{})
	if !ok || len(listMap)<=0 {
		log.Println("busiHeartbeat getVersion no list")
		return nil
	}

	row,ok:=listMap[0].(map[string]interface{})
	if !ok {
		log.Println("busiHeartbeat getVersion convert list[0] to map error")
		return nil
	}

	return row[crv.CC_VERSION]
}

func (busi *Busi)DealDeviceHeartbeat(deviceID,vin string){
	//登录
	//if busi.CrvClient.Login() ==0 {
		//添加心跳记录到记录表
		saveReq:=getDeviceHeartbeatSaveReq(deviceID,vin)
		busi.CrvClient.Save(saveReq,"")

		//查询车辆信息
		queryReq:=getQueryVinReq(vin)
		rsp,err:=busi.CrvClient.Query(queryReq,"")
		if err==0 && rsp.Error==false {
			version:=getVersion(rsp.Result)
			if version != nil {
				//更新车辆对应设备信息
				saveReq=getUpdateVinDeviceReq(deviceID,vin,version)
				busi.CrvClient.Save(saveReq,"")
			}
		} else {
			log.Printf("Query vin info error %s,%d",rsp.Message,err)
		}
	//}
}
