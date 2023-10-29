package send

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
)

var QueryVehicleFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "DeviceNumber"},
	{"field": "ProjectNum"},
	{"field": "TestSpecification"},
}

type sendVehicleItem struct {
	Vin string
	DeviceID string
	ProjectNum string
	TestSpecification string
}

type sendVehicle struct {
	CRVClient *crv.CRVClient
}

func (sv *sendVehicle)getVehicleIDs(row map[string]interface{})([]string){
	cars,_:=row["cars"].(map[string]interface{})["list"].([]interface{})
	ids:=[]string{}
	for _,car:=range(cars){
		carMap,_:=car.(map[string]interface{})
		id,_:=carMap["id"]
		strID,_:=id.(string)
		ids=append(ids,strID)
	}
	return ids
}

func (sv *sendVehicle)getSendVehicle(row map[string]interface{},token string)([]sendVehicleItem,int){
	carIDs:=sv.getVehicleIDs(row)
	rsp,errorCode:=sv.queryVehicleList(carIDs,token)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	return sv.convertVehicleList(rsp)
}

func (sv *sendVehicle)queryVehicleList(carIDs []string,token string)(*crv.CommonRsp,int){
	log.Println("start queryDtcList")
	commonRep:=crv.CommonReq{
		ModelID:"vehiclemanagement",
		Filter:&map[string]interface{}{
			"id":&map[string]interface{}{
				"Op.in":carIDs,
			},
		},
		Fields:&QueryVehicleFields,
	}

	return sv.CRVClient.Query(&commonRep,token)
} 

func (sv *sendVehicle)convertVehicleList(queryResult *crv.CommonRsp)([]sendVehicleItem,int){
	log.Println("start convertVehicleList")
	list,ok:=queryResult.Result["list"]
	if !ok {
		log.Println("convertVehicleList queryResult no list")
		return nil,common.ResultNoNoVehicle
	}

	vcList,ok:=list.([]interface{})
	if !ok || len(vcList)<=0 {
		log.Println("convertVehicleList queryResult no list")
		return nil,common.ResultNoNoVehicle
	}

	sendVehicles:=make([]sendVehicleItem,len(vcList))

	for index,item:=range vcList {
		mapItem:=item.(map[string]interface{})
		sendVehicles[index].Vin=mapItem["id"].(string)
		sendVehicles[index].DeviceID=mapItem["DeviceNumber"].(string)
		sendVehicles[index].ProjectNum=mapItem["ProjectNum"].(string)
		sendVehicles[index].TestSpecification=mapItem["TestSpecification"].(string)
	}	

	return sendVehicles,common.ResultSuccess
}

