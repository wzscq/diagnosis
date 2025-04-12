package fulldata

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"errors"
)

type FullDataConf struct {
    ID string
	Name string
}

var queryConfFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "name"},
}

func GetConf(token,id string,crvClient *crv.CRVClient)(*FullDataConf,error){
	commonRep := crv.CommonReq{
		ModelID: "full_data_collection_conf",
		Filter: &map[string]interface{}{
			"id": id,
		},
		Fields: &queryConfFields,
	}

	rsp, errCode := crvClient.Query(&commonRep, token)
	if errCode != common.ResultSuccess {
		return nil, errors.New("获取配置参数失败")
	}

	list, ok := rsp.Result["list"].([]interface{})
	if !ok || len(list) == 0 {
		return nil, errors.New("获取配置参数失败")
	}

	row, ok := list[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("获取配置参数失败")
	}

	conf := &FullDataConf{
		ID:   row["id"].(string),
		Name: row["name"].(string),
	}

	return conf, nil
}

type CarInfo struct {
    ID string
    VehicleManagementCode string
    DeviceNumber string
	ProjectNum   string
}

var queryCarFields = []map[string]interface{}{
    {"field": "id"},
	{"field":"VehicleManagementCode"},
    {"field": "DeviceNumber"},
    {"field": "ProjectNum"},
}

func GetCar(token, carID string, crvClient *crv.CRVClient) (*CarInfo, error) {
    commonRep := crv.CommonReq{
        ModelID: "view_vehicle",
        Filter: &map[string]interface{}{
            "id": carID,
        },
        Fields: &queryCarFields,
    }

    rsp, errCode := crvClient.Query(&commonRep, token)
    if errCode != common.ResultSuccess {
        return nil, errors.New("获取车辆信息失败")
    }

    list, ok := rsp.Result["list"].([]interface{})
    if !ok || len(list) == 0 {
        return nil, errors.New("获取车辆信息失败")
    }

    row, ok := list[0].(map[string]interface{})
    if !ok {
        return nil, errors.New("获取车辆信息失败")
    }

    car := &CarInfo{
        ID:       row["id"].(string),
		VehicleManagementCode:row["VehicleManagementCode"].(string),
        DeviceNumber: row["DeviceNumber"].(string),
        ProjectNum: row["ProjectNum"].(string),
    }

    return car, nil
}

func CreateSendRec(token string,carInfo *CarInfo,conf *FullDataConf,crvClient *crv.CRVClient)(error){
	commonRep:=crv.CommonReq{
		ModelID:"full_data_conf_send_rec",
		List:&[]map[string]interface{}{
			{
				"_save_type":"create",
				"vehicle_code":carInfo.VehicleManagementCode,
				"project_num":carInfo.ProjectNum,
				"device_num":carInfo.DeviceNumber,
				"collection_conf_id":conf.ID,
				"status":"1",
			},
		},
	}

	_,errCode:=crvClient.Save(&commonRep,token)
	if errCode!=common.ResultSuccess {
		return errors.New("创建发送记录失败")
	}

	return nil
}