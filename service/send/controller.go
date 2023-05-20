package send

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
	"digimatrix.com/diagnosis/common"
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/mqtt"
	"bytes"
)

type SendController struct {
	CRVClient *crv.CRVClient
	MQTTClient *mqtt.MQTTClient
	SendRecordCache *SendRecordCache
	DeviceSignalCache *DeviceSignalCache
	FilePath string
	DBCUploadTopic string
}

func (controller *SendController) sendParameter1 (c *gin.Context){
	log.Println("start sendParameter")

	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end redirect with error")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
    }	

	if rep.List==nil || len(*rep.List)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//controller.CRVClient.Token=header.Token
	//生成下发参数
	/*errorCode:=controller.CRVClient.Login()
	if errorCode!=0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultCannotLoginCRV,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}*/

	sp:=&SendParameter{
		CRVClient:controller.CRVClient,
		SignalList:&map[string]interface{}{},
	}
	parameter,errorCode:=sp.getSendParameter((*rep.List)[0],header.Token)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	log.Println(parameter)
	bf := bytes.NewBuffer([]byte{})
    jsonEncoder := json.NewEncoder(bf)
    jsonEncoder.SetEscapeHTML(false)
    jsonEncoder.Encode(parameter)
	strParam:=bf.String()

	log.Println(strParam)
	rsp1:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp1)
	return

	//获取下发车辆列表
	sv:=&sendVehicle{
		CRVClient:controller.CRVClient,
	}
	vehicles,errorCode:=sv.getSendVehicle((*rep.List)[0],header.Token)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	log.Println(vehicles)

	//创建下发记录
	_,errorCode=createSendRecords(
		controller.CRVClient,
		controller.SendRecordCache,
		vehicles,
		rep.UserID,
		strParam,
		header.Token)

	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	//log.Println(*saveRsp)
	//执行参数下发
	errorCode=sendByMqtt(controller.MQTTClient,vehicles,strParam)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	errorCode=sendSignalByMqtt(
		controller.MQTTClient,
		vehicles,
		sp.SignalList,
		controller.DeviceSignalCache,
		"diag",
		controller.CRVClient,
		controller.SendRecordCache,
		rep.UserID,
		header.Token)

	rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("end sendParameter")
}

func (controller *SendController)uploadDBC(c *gin.Context){
	log.Println("start uploadDBC")
	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
    }	

	if rep.List==nil || len(*rep.List)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	log.Println(rep)
	
	//保存文件
	row:=(*rep.List)[0]
	fieldVal,_:=row["dbc_file"]
	fileName,errorCode:=saveDBCFile(controller.FilePath,fieldVal)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//发送消息到后台解析程序
	//构造发送结构
	content:=map[string]interface{}{
		"user_id":rep.UserID,
		"platform_id":row["platform_id"],
		"domain_id":row["domain_id"],
		"dbc_file":fileName,
	}

	jsonBytes,_:=json.Marshal(content)
	log.Println(string(jsonBytes))
	errorCode=controller.MQTTClient.Publish(controller.DBCUploadTopic,string(jsonBytes))
	
	rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (controller *SendController)sendEventParameter (c *gin.Context){
	log.Println("start sendEventParameter")
	
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end redirect with error")
		return
	}	
	
	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
    }	

	if rep.List==nil || len(*rep.List)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//controller.CRVClient.Token=header.Token
	//生成下发参数
	/*errorCode:=controller.CRVClient.Login()
	if errorCode!=0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultCannotLoginCRV,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}*/

	sp:=&SendEventParameter{
		CRVClient:controller.CRVClient,
		SignalList:&map[string]interface{}{},
	}
	parameter,errorCode:=sp.getSendParameter((*rep.List)[0],header.Token)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	log.Println(parameter)
	bf := bytes.NewBuffer([]byte{})
    jsonEncoder := json.NewEncoder(bf)
    jsonEncoder.SetEscapeHTML(false)
    jsonEncoder.Encode(parameter)
	strParam:=bf.String()
	log.Println(strParam)
	//获取下发车辆列表
	sv:=&sendVehicle{
		CRVClient:controller.CRVClient,
	}
	vehicles,errorCode:=sv.getSendVehicle((*rep.List)[0],header.Token)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	log.Println(vehicles)

	//创建下发记录
	_,errorCode=createSendEventRecords(
		controller.CRVClient,
		controller.SendRecordCache,
		vehicles,
		rep.UserID,
		strParam,
		header.Token)

	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	//log.Println(*saveRsp)
	//执行参数下发
	errorCode=sendEventByMqtt(controller.MQTTClient,vehicles,strParam)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	errorCode=sendSignalByMqtt(
		controller.MQTTClient,
		vehicles,
		sp.SignalList,
		controller.DeviceSignalCache,
		"event",
		controller.CRVClient,
		controller.SendRecordCache,
		rep.UserID,
		header.Token)	

	rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("end sendEventParameter")
}

//Bind bind the controller function to url
func (controller *SendController) Bind(router *gin.Engine) {
	log.Println("Bind controller")
	router.POST("/sendParameter1", controller.sendParameter1)
	router.POST("/uploadDBC",controller.uploadDBC)
	router.POST("/sendEventParameter",controller.sendEventParameter)
}