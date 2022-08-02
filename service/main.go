package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"digimatrix.com/diagnosis/dashboard"
	"digimatrix.com/diagnosis/common"
	"digimatrix.com/diagnosis/report"
	"digimatrix.com/diagnosis/send"
	"digimatrix.com/diagnosis/mqtt"
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/busi"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowAllOrigins:true,
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
        AllowCredentials: true,
    }))

	conf:=common.InitConfig()

	//crvClinet 用于到crvframeserver的请求
	crvClinet:=crv.CRVClient{
		Server:conf.CRV.Server,
		User:conf.CRV.User,
		Password:conf.CRV.Password,
		AppID:conf.CRV.AppID,
	}

	//实际的业务处理模块
	busiModule:=busi.Busi{
		CrvClient:&crvClinet,
	}

	//mqttclient
	mqttClient:=mqtt.MQTTClient{
		Broker:conf.MQTT.Broker,
		User:conf.MQTT.User,
		Password:conf.MQTT.Password,
		HeartbeatTopic:conf.MQTT.HeartbeatTopic,
		Busi:&busiModule,
	}
	mqttClient.Init()
	
	repo:=&dashboard.DefatultRepository{}
    repo.Connect(
        conf.Mysql.Server,
        conf.Mysql.User,
        conf.Mysql.Password,
        conf.Mysql.DBName)
	dashboardController:=&dashboard.Controller{
    	Repository:repo,
    }
    dashboardController.Bind(router)

	reportController:=report.CreateController(
		conf.Mongo.Server, 
		conf.Mongo.DBName,
		conf.Mongo.User,
        conf.Mongo.Password)
	reportController.Bind(router)

	sendController:=send.SendController{}
	sendController.Bind(router)

	router.Run(conf.Service.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}