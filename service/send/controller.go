package send

import (
	"log"
	"encoding/json"
	"bytes"
	"fmt"
	"net/http"
	"crypto/md5"
	"strconv"
	"github.com/gin-gonic/gin"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type loginRep struct {
    UserID     string  `json:"userID"`
    Password  string   `json:"password"`
	AppID     string   `json:"appID"`
} 

type loginResult struct {
    UserID     *string  `json:"userID"`
    UserName  *string  `json:"userName"`
	Token     *string  `json:"token"`
	AppID     *string  `json:"appID"`
}

type LoginRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result loginResult `json:"result"`
}

type CommonRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result map[string]interface{} `json:"result"`
}

type commonRep struct {
	ModelID string `json:"modelID"`
	ViewID *string `json:"viewID"`
	Filter *map[string]interface{} `json:"filter"`
	List *[]map[string]interface{} `json:"list"`
	Fields *[]map[string]interface{} `json:"fields"`
	//Sorter *[]sorter `json:"sorter"`
	//SelectedRowKeys *[]string `json:"selectedRowKeys"`
	//Pagination *pagination `json:"pagination"`
}

type SendController struct {
}

var QueryConfigFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "state_value"},
	{"field": "wait_time"},
	{
		"field": "ecu_state_machine",
		"fieldType": "many2one",
		"relatedModelID": "diag_signal",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "name"},
			{"field": "can_id"},
			{"field": "byte_order"},
			{"field": "start_addr"},
			{"field": "len"},
			{"field": "factor"},
			{"field": "offset"},
			{"field": "pdu_id"},
		},
	},
	{
		"field": "mileage_name",
		"fieldType": "many2one",
		"relatedModelID": "diag_signal",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "name"},
			{"field": "can_id"},
			{"field": "byte_order"},
			{"field": "start_addr"},
			{"field": "len"},
			{"field": "factor"},
			{"field": "offset"},
			{"field": "pdu_id"},
		},
	},
}

var QueryDTCSignalFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "dtc_id"},
	{
		"field": "signals",
		"fieldType": "many2many",
		"relatedModelID": "diag_signal",
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		 	{"field": "name"},
			{"field": "can_id"},
			{"field": "start_addr"},
			{"field": "pdu_id"},
			{"field": "factor"},
			{"field": "offset"},
		},
	},
}

var QueryParameterFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "version"},
	{"field": "name"},
	{"field": "platform_id"},
	{"field": "time_offset"},
	{"field": "use_triggercanid"},
	/*{
		"field": "domain_id",
		"fieldType": "MANY_TO_ONE",
		"relatedModelID": "diag_domain",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "name"},
		},
	},*/
	{
		"field": "ecu_id",
		"fieldType": "many2one",
		"relatedModelID": "diag_ecu",
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		 	{"field": "name"},
			{"field": "tx"},
			{"field": "rx"},
			{"field": "instruct"},
			{"field": "trigger_can_id"},
		},
	},
	{
		"field": "logistics_did",
		"fieldType": "many2many",
		"relatedModelID": "diag_logistics",
		"associationModelID": "diag_parameter_logistics_did",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "did"},
		},
	},
	{
		"field": "internal_did",
		"fieldType": "many2many",
		"relatedModelID": "diag_logistics",
		"associationModelID": "diag_parameter_internal_did",
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		  	{"field": "did"},
		},
	},
	{
		"field": "triggers",
		"fieldType": "one2many",
		"relatedModelID": "diag_parameter_trigger",
		"relatedField": "parameter_id",
		"fields": []map[string]interface{}{
		  	{"field": "id"},
		  	{
				"field": "diag_signal_id",
				"fieldType": "many2one",
				"relatedModelID": "diag_signal",
				"fields": []map[string]interface{}{
			  		{"field": "id"},
			  		{"field": "name"},
					{"field": "can_id"},
					{"field": "start_addr"},
					{"field": "pdu_id"},
					{"field": "factor"},
					{"field": "offset"},
				},
		  	},
		  	{"field": "logic"},
		  	{"field": "value"},
		  	//{"field": "version"},
		  	{"field": "parameter_id"},
		},
	},
	//{"field": "remark"},
}

func (controller *SendController) login()(string,int) {
	log.Println("start login")
	loginRep:=loginRep{
		UserID:"carapi",
		Password:"a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
		AppID:"diagnosis",
	}

	postJson,_:=json.Marshal(loginRep)
	postBody:=bytes.NewBuffer(postJson)
	resp,err:=http.Post("http://127.0.0.1:8200/user/login","application/json",postBody)
	log.Println("start login2")
	if err != nil {
		log.Println("login error",err)
		return "",-1
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 { 
		log.Println("login error",resp)
		return "",-1
	}

	decoder := json.NewDecoder(resp.Body)
	loginRsp:=LoginRsp{}
	err = decoder.Decode(&loginRsp)
	if err != nil {
		log.Println("result decode failed [Err:%s]", err.Error())
		return "",-1
	}

	if loginRsp.Result.Token == nil {
		log.Println("login error ",loginRsp)
		return "",-1
	}

	log.Println("login success")
	return *loginRsp.Result.Token,0
}

func (controller *SendController)getParameterID(rep commonRep)(string,int){
	if rep.List == nil {
		return "",-1
	}

	if len(*rep.List) <=0 {
		return "",-1
	}

	idObj,ok:=(*rep.List)[0]["id"]
	if !ok {
		return "",-1
	}

	idStr,ok:=idObj.(string)
	if !ok {
		return "",-1
	}

	return idStr,0
}

func (controller *SendController)queryData(commonRep commonRep,token string)(*CommonRsp,int){
	log.Println("start queryData ...")
	postJson,_:=json.Marshal(commonRep)
	postBody:=bytes.NewBuffer(postJson)
	req,err:=http.NewRequest("POST", "http://127.0.0.1:8200/data/query",postBody)
	if err != nil {
		log.Println("queryData NewRequest error",err)
		return nil,-1
	}
	req.Header.Set("token", token)
	req.Header.Set("Content-Type","application/json")
	
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("queryData Do request error",err)
		return nil,-1
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 { 
		log.Println("queryData StatusCode error",resp)
		return nil,-1
	}

	decoder := json.NewDecoder(resp.Body)
	commonRsp:=CommonRsp{}
	err = decoder.Decode(&commonRsp)
	if err != nil {
		log.Println("queryData result decode failed [Err:%s]", err.Error())
		return nil,-1
	}

	resultJson,_:=json.Marshal(&commonRsp.Result)
	log.Println(string(resultJson))

	log.Println("end queryData success")
	return &commonRsp,0
}

func (controller *SendController)queryDiagParameter(token,parameterID string)(*CommonRsp,int){
	log.Println("start queryDiagParameter")
	commonRep:=commonRep{
		ModelID:"diag_parameter",
		Filter:&map[string]interface{}{
			"id":parameterID,
		},
		Fields:&QueryParameterFields,
	}

	return controller.queryData(commonRep,token)
}

func (controller *SendController)queryDtcList(token string,ecuID string)(*CommonRsp,int){
	log.Println("start queryDtcList")
	commonRep:=commonRep{
		ModelID:"diag_manual_fault",
		Filter:&map[string]interface{}{
			"ecu_id":ecuID,
		},
		Fields:&QueryDTCSignalFields,
	}

	return controller.queryData(commonRep,token)
}

func (controller *SendController)queryDiagConfig(token string)(*CommonRsp,int){
	log.Println("start queryDiagConfig")
	commonRep:=commonRep{
		ModelID:"diag_exe_config",
		Fields:&QueryConfigFields,
	}

	return controller.queryData(commonRep,token)
}

func (controller *SendController)convertDiagConf(queryResult *CommonRsp)(*map[string]interface{},int){
	log.Println("start convertDiagConf")
	if queryResult.Error {
		log.Println("convertDiagConf queryResult has error")
		return nil,-1
	}

	list,ok:=queryResult.Result["list"]
	if !ok {
		log.Println("convertDiagConf queryResult no list")
		return nil,-1
	}

	listMap,ok:=list.([]interface{})
	if !ok || len(listMap)<=0 {
		log.Println("convertDiagConf queryResult no list")
		return nil,-1
	}

	row,ok:=listMap[0].(map[string]interface{})
	if !ok {
		log.Println("convertDiagConf convert list[0] to map error")
		return nil,-1
	}

	conf:=map[string]interface{}{}
	//conf["TriggerCanId"]=row["trigger_can_id"]
	ecuMachineLst,ok:=row["ecu_state_machine"].(map[string]interface{})["list"].([]interface{})
	if ok && len(ecuMachineLst)>0 {
		ecuMachine,_:=ecuMachineLst[0].(map[string]interface{})
		executionPara:=map[string]interface{}{}
		executionPara["StateValue"]=row["state_value"]
		executionPara["WaitTime"]=row["wait_time"]
		executionPara["EcuStateMachine"]=ecuMachine["id"]
		
		executionPara["CanID"]=ecuMachine["can_id"]
		executionPara["ByteOrder"]=ecuMachine["byte_order"]
		executionPara["StartAddr"]=ecuMachine["start_addr"]
		executionPara["Len"]=ecuMachine["len"]
		executionPara["Factor"]=ecuMachine["factor"]
		executionPara["Offset"]=ecuMachine["offset"]
		executionPara["PduID"]=ecuMachine["pdu_id"]
		conf["ExecutionPara"]=executionPara
	}
	
	mileageLst,ok:=row["mileage_name"].(map[string]interface{})["list"].([]interface{})
	if ok && len(mileageLst)>0 {
		mileage,_:=mileageLst[0].(map[string]interface{})
		mileagePara:=map[string]interface{}{}
		mileagePara["MileageName"]=mileage["id"]

		mileagePara["CanID"]=mileage["can_id"]
		mileagePara["ByteOrder"]=mileage["byte_order"]
		mileagePara["StartAddr"]=mileage["start_addr"]
		mileagePara["Len"]=mileage["len"]
		mileagePara["Factor"]=mileage["factor"]
		mileagePara["Offset"]=mileage["offset"]
		mileagePara["PduID"]=mileage["pdu_id"]
		conf["MileagePara"]=mileagePara
	}

	return &conf,0 
}

func (controller *SendController)convertDtcList(queryResult *CommonRsp)(*[]interface{},int){
	log.Println("start convertDtcList")
	if queryResult.Error {
		log.Println("convertDtcList queryResult has error")
		return nil,-1
	}

	list,ok:=queryResult.Result["list"]
	if !ok {
		log.Println("convertDtcList queryResult no list")
		return nil,-1
	}

	dtcList,ok:=list.([]interface{})
	if !ok || len(dtcList)<=0 {
		log.Println("convertDtcList queryResult no list")
		return nil,-1
	}

	for _,item:=range dtcList {
		mapItem:=item.(map[string]interface{})
		mapItem["DtcNum"]=mapItem["dtc_id"]
		signals,_:=mapItem["signals"].(map[string]interface{})["list"].([]interface{})
		for _,signal:=range signals {
			signalItem:=signal.(map[string]interface{})
			sCanID,_:=signalItem["can_id"].(string)
			sPduID,_:=signalItem["pdu_id"].(string)
			sStartAddr,_:=signalItem["start_addr"].(string)
			mapItem["SignalID"]=fmt.Sprintf("%s:%s:%s",sCanID,sPduID,sStartAddr)
			//signalItem["SignalID"]=signalItem["id"]
			signalItem["SignalName"]=signalItem["name"]
			delete(signalItem,"id")
			delete(signalItem,"name")
			delete(signalItem,"can_id")
			delete(signalItem,"pdu_id")
			delete(signalItem,"start_addr")
		}
		mapItem["CorrelationSignal"]=signals
		delete(mapItem,"id")
		delete(mapItem,"dtc_id")
		delete(mapItem,"signals")
	}	

	return &dtcList,0
}

func (controller *SendController)convertDiagParameter(qeuryResult *CommonRsp)(*map[string]interface{},int){
	log.Println("start convertDiagParameter")
	if qeuryResult.Error {
		log.Println("convertDiagParameter qeuryResult has error")
		return nil,-1
	}

	list,ok:=qeuryResult.Result["list"]
	if !ok {
		log.Println("convertDiagParameter qeuryResult no list")
		return nil,-1
	}

	listMap,ok:=list.([]interface{})
	if !ok || len(listMap)<=0 {
		log.Println("convertDiagParameter qeuryResult no list")
		return nil,-1
	}

	row,ok:=listMap[0].(map[string]interface{})
	if !ok {
		log.Println("convertDiagParameter convert list[0] to map error")
		return nil,-1
	}

	diagParaList:=map[string]interface{}{}
	diagParaList["TimeOffset"]=row["time_offset"]
	diagParaList["Channel"]=row["use_triggercanid"]
	
	ecuList,ok:=row["ecu_id"].(map[string]interface{})["list"].([]interface{})
	if ok && len(ecuList)>0 {
		ecu,_:=ecuList[0].(map[string]interface{})
		diagParaList["Ecu"]=ecu["id"]
		diagParaList["TxId"]=ecu["tx"]
		diagParaList["RxId"]=ecu["rx"]
		diagParaList["TriggerCanId"]=ecu["trigger_can_id"]
		diagParaList["DiagInstruct"]=ecu["instruct"]
	}
	
	logistics,_:=row["logistics_did"].(map[string]interface{})["list"].([]interface{})
	for _,item:=range logistics {
		mapItem:=item.(map[string]interface{})
		delete(mapItem,"id")
		mapItem["DID"]=mapItem["did"]
		delete(mapItem,"did")
	}

	internallogistics,_:=row["internal_did"].(map[string]interface{})["list"].([]interface{})
	for _,item:=range internallogistics {
		mapItem:=item.(map[string]interface{})
		delete(mapItem,"id")
		mapItem["DID"]=mapItem["did"]
		delete(mapItem,"did")
	}

	triggerList,_:=row["triggers"].(map[string]interface{})["list"].([]interface{})
	for _,item:=range triggerList {
		mapItem:=item.(map[string]interface{})
		delete(mapItem,"id")
		delete(mapItem,"parameter_id")
		signalList,ok:=mapItem["diag_signal_id"].(map[string]interface{})["list"].([]interface{})
		var floatFactor,floatOffset float64
		var err error
		if ok && len(signalList)>0 {
			signal,_:=signalList[0].(map[string]interface{})
			mapItem["SignalName"]=signal["name"]
			
			sCanID,_:=signal["can_id"].(string)
			sPduID,_:=signal["pdu_id"].(string)
			sStartAddr,_:=signal["start_addr"].(string)
			mapItem["SignalID"]=fmt.Sprintf("%s:%s:%s",sCanID,sPduID,sStartAddr)

			sFactor,_:=signal["factor"].(string)
			floatFactor, err = strconv.ParseFloat(sFactor, 64)
			if err!=nil {
				log.Println("trigger signal factor can not convert to float64")
				log.Println("factor value is :"+sFactor)
				return nil,-1
			}
			sOffset,_:=signal["offset"].(string)
			floatOffset, err= strconv.ParseFloat(sOffset, 64)
			if err!=nil {
				log.Println("trigger signal offset can not convert to float64")
				return nil,-1
			}
		}
		delete(mapItem,"diag_signal_id")
		mapItem["Logic"]=mapItem["logic"]
		delete(mapItem,"logic")
		sValue,_:=mapItem["value"].(string)
		floatVal, err := strconv.ParseFloat(sValue, 64)
		if err!=nil {
			log.Println("trigger value can not convert to float64")
			return nil,-1
		}
		mapItem["Value"]=floatVal/floatFactor+floatOffset
		delete(mapItem,"value")
	}
	
	diagParaList["Logistics"]=logistics
	diagParaList["Internallogistics"]=internallogistics
	diagParaList["TriggerList"]=triggerList

	return &diagParaList,0
}

func (controller *SendController) connectHandler(client mqtt.Client){
	log.Println("connectHandler connect status: ",client.IsConnected())
}

func (controller *SendController) connectLostHandler(client mqtt.Client, err error){
	log.Println("connectLostHandler connect status: ",client.IsConnected(),err)
}

func (controller *SendController) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("messagePublishHandler topic: ",msg.Topic())
}

func (controller *SendController) getMqttClient()(*mqtt.Client){
	broker := "49.4.3.226"
	port := 1983
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d",broker,port))
	opts.SetClientID("diagonsis_mqtt_client")
	opts.SetUsername("mosquitto")
	opts.SetPassword("123456")
	opts.SetDefaultPublishHandler(controller.messagePublishHandler)
	opts.OnConnect = controller.connectHandler
	opts.OnConnectionLost = controller.connectLostHandler
	client:=mqtt.NewClient(opts)
	if token:=client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error)
		return nil
	}
	return &client
}

func (controller *SendController) mergPara(
	diagParaList,diagConf *map[string]interface{},
	dtcList *[]interface{},
	projectNum interface{})(*map[string]interface{},int){
	log.Println("start mergPara")
	distPara:=*diagConf
	distPara["DiagParaList"]=diagParaList
	distPara["DtcList"]=dtcList

	postJson,_:=json.Marshal(distPara)
	log.Println("postJson: ",string(postJson))

	md5Inst := md5.New()
	md5Inst.Write([]byte(string(postJson)))
	md5str := md5Inst.Sum([]byte(""))

	distPara["MD5"]=md5str
	distPara["ProjectNum"]=projectNum

	return &distPara,0
}

func (controller *SendController) distributeToMqtt(content *map[string]interface{})(int){
	log.Println("start distributeToMqtt")
	client:=controller.getMqttClient()
	if client == nil {
		return -1
	}
	defer (*client).Disconnect(250)

	postJson,_:=json.Marshal(content)
	log.Println("content: ",string(postJson))
	token:=(*client).Publish("DiagConfUpload",0,false,string(postJson))
	token.Wait()

	log.Println("end distributeToMqtt")
	return 0
}

func (controller *SendController)getProjectNum(qeuryResult *CommonRsp)(interface{},int){
	log.Println("start convertDiagParameter")
	if qeuryResult.Error {
		log.Println("convertDiagParameter qeuryResult has error")
		return nil,-1
	}

	list,ok:=qeuryResult.Result["list"]
	if !ok {
		log.Println("convertDiagParameter qeuryResult no list")
		return nil,-1
	}

	listMap,ok:=list.([]interface{})
	if !ok || len(listMap)<=0 {
		log.Println("convertDiagParameter qeuryResult no list")
		return nil,-1
	}

	row,ok:=listMap[0].(map[string]interface{})
	if !ok {
		log.Println("convertDiagParameter convert list[0] to map error")
		return nil,-1
	}

	return row["platform_id"],0
}

func (controller *SendController) sendParameter (c *gin.Context){
	log.Println("start sendParameter")
	var rep commonRep
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "")
		return
    }	
	log.Println(rep)
	//获取配置记录ID
	parameterID,errorCode:=controller.getParameterID(rep)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	//登录平台查询数据
	token,errorCode:=controller.login()
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	//获取参数配置
	queryResult,errorCode:=controller.queryDiagParameter(token,parameterID)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	//ProjectNum
	projectNum,errorCode:=controller.getProjectNum(queryResult)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	//解析成待发送格式
	diagParaList,errorCode:=controller.convertDiagParameter(queryResult)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	//获取DTC对应的信号列表
	queryResult,errorCode=controller.queryDtcList(token,(*diagParaList)["Ecu"].(string))
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	//解析成待发送格式
	dtcList,errorCode:=controller.convertDtcList(queryResult)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	//获取诊断参数配置
	queryResult,errorCode=controller.queryDiagConfig(token)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	//解析成待发送格式
	diagConf,errorCode:=controller.convertDiagConf(queryResult)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	//组装发送json
	distJson,errorCode:=controller.mergPara(diagParaList,diagConf,dtcList,projectNum)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	errorCode=controller.distributeToMqtt(distJson)
	if errorCode !=0 {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, "")
	log.Println("end sendParameter")
}

//Bind bind the controller function to url
func (controller *SendController) Bind(router *gin.Engine) {
	log.Println("Bind controller")
	router.POST("/sendParameter", controller.sendParameter)
}