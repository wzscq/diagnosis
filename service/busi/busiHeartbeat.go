package busi

import (
	"log"
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"time"
)

const (
	MODEL_HEARTBEAT = "diag_device_heartbeat"
	MODEL_VEHICLE  = "vehiclemanagement"
)

func getDeviceHeartbeatSaveReq(deviceID,vin string,veicheleRow map[string]interface{})(*crv.CommonReq){
	if veicheleRow == nil {
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

	return &crv.CommonReq{
		ModelID:MODEL_HEARTBEAT,
		List:&[]map[string]interface{}{
			map[string]interface{}{
				"device_number":deviceID,
				"vin":vin,
				"vehicle_management_code":veicheleRow["VehicleManagementCode"],
				"project_num":veicheleRow["ProjectNum"],
				"test_specification":veicheleRow["TestSpecification"],
				"develop_phase":veicheleRow["developPhase"],
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
	now:=time.Now().Format("2006-01-02 15:04:05")

	return &crv.CommonReq{
		ModelID:MODEL_VEHICLE,
		Fields:&[]map[string]interface{}{
			{"field": "id"},
			{"field": "version"},
			{"field": "VehicleManagementCode"},
			{"field": "ProjectNum"},
			{"field": "TestSpecification"},
			{"field": "developPhase"},
		},
		Filter:&map[string]interface{}{
			"id":vin,
			"BindingDate":map[string]interface{}{
				"Op.lte":now,
			},
			"UntieDate":map[string]interface{}{
				"Op.gte":now,
			},
		},
		Pagination:&crv.Pagination{
			Current:1,
			PageSize:1,
		},
	}
}

func getVeichleRow(result map[string]interface{})(map[string]interface{}){
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

	return row
}

func (busi *Busi)DealDeviceHeartbeat(deviceID,vin string){
	//登录
	//if busi.CrvClient.Login() ==0 {
		//查询车辆信息
		queryReq:=getQueryVinReq(vin)
		rsp,err:=busi.CrvClient.Query(queryReq,"")
		if err==common.ResultSuccess && rsp.Error==false {
			//添加心跳记录到记录表
			veicheleRow:=getVeichleRow(rsp.Result)
			saveReq:=getDeviceHeartbeatSaveReq(deviceID,vin,veicheleRow)
			busi.CrvClient.Save(saveReq,"")
		} else {
			log.Printf("Query vin info error %s,%d",rsp.Message,err)
		}
	//}
}
