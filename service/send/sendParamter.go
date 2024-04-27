package send


import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
	"crypto/md5"
	"encoding/json"
)

type SendParameter struct {
	CRVClient *crv.CRVClient
	SignalList *map[string]interface{}
}

func (sp *SendParameter)getParameterIDs(row map[string]interface{})([]string){
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

func (sp *SendParameter)generateSendObject(
	dp *diagParameter,
	dtc *dtcList,
	dc *diagConf)(map[string]interface{}){

	log.Println("start generateSendObject")
	distPara:=dc.Conf
	distPara["DiagParaList"]=dp.Params
	distPara["DtcList"]=dtc.DtcList
	distPara["ProjectNum"]=dp.getProjectNum()

	postJson,_:=json.Marshal(distPara)
	log.Println("postJson: ",string(postJson))

	md5Inst := md5.New()
	md5Inst.Write([]byte(string(postJson)))
	md5str := md5Inst.Sum([]byte(""))

	distPara["MD5"]=md5str
	
	return distPara
}

func (sp *SendParameter)getSendParameter(row map[string]interface{},token string)(map[string]interface{},int){
	log.Println("getSendParameter start")
	ids:=sp.getParameterIDs(row)
	//获取参数记录
	log.Println("getSendParameter getDiagParams ...")
	dp,errorCode:=getDiagParams(ids,sp.CRVClient,sp.SignalList,token)
	if errorCode!=common.ResultSuccess {
		log.Println("getSendParameter getDiagParams error")
		return nil,errorCode
	}

	//获取DTC
	log.Println("getSendParameter getDtcList ...")
	dtc,errorCode:=getDtcList(dp.PlatformID,dp.getEcuIDs(),dp.getEcuChannelMap(),sp.CRVClient,sp.SignalList,token)
	if errorCode!=common.ResultSuccess {
		log.Println("getSendParameter getDtcList error")
		return nil,errorCode
	}
	//获取配置信息
	log.Println("getSendParameter getDiagConf ...")
	dc,errorCode:=getDiagConf(sp.CRVClient,dp.PlatformID,token)
	if errorCode!=common.ResultSuccess {
		log.Println("getSendParameter getDiagConf error")
		return nil,errorCode
	}
	//组装发送JSON
	result:=sp.generateSendObject(dp,dtc,dc)

	return result,common.ResultSuccess
}