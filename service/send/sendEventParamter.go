package send


import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
	"crypto/md5"
	"encoding/json"
)

type SendEventParameter struct {
	CRVClient *crv.CRVClient
	SignalList *map[string]interface{}
}

func (sp *SendEventParameter)getParameterIDs(row map[string]interface{})([]string){
	paramList,_:=row["parameters"].(map[string]interface{})["list"].([]interface{})
	ids:=[]string{}
	for _,paramRow:=range(paramList){
		rowMap,_:=paramRow.(map[string]interface{})
		id,_:=rowMap["id"]
		strID,_:=id.(string)
		ids=append(ids,strID)
	}
	return ids
}

func (sp *SendEventParameter)generateSendObject(
	dp *eventParameter)(map[string]interface{}){

	log.Println("start generateSendObject")
	distPara:=map[string]interface{}{}
	distPara["EventList"]=dp.Params
	
	postJson,_:=json.Marshal(distPara)
	log.Println("postJson: ",string(postJson))

	md5Inst := md5.New()
	md5Inst.Write([]byte(string(postJson)))
	md5str := md5Inst.Sum([]byte(""))

	distPara["MD5"]=md5str
	
	return distPara
}

func (sp *SendEventParameter)getSendParameter(row map[string]interface{})(map[string]interface{},int){
	log.Println("getSendParameter start")
	ids:=sp.getParameterIDs(row)
	//获取参数记录
	log.Println("getSendParameter getEventParams ...")
	dp,errorCode:=getEventParams(ids,sp.CRVClient,sp.SignalList)
	if errorCode!=common.ResultSuccess {
		log.Println("getSendParameter getEventParams error")
		return nil,errorCode
	}
	//组装发送JSON
	result:=sp.generateSendObject(dp)

	return result,common.ResultSuccess
}