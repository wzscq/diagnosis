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
	"digimatrix.com/diagnosis/oauth"
	"digimatrix.com/diagnosis/idm"
	"digimatrix.com/diagnosis/saicinterface"
	"log"
	"time"
)

func main() {
	//设置log打印文件名和行号
  log.SetFlags(log.Lshortfile | log.LstdFlags)

	//初始化时区
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

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
		Token:conf.CRV.Token,
		AppID:conf.CRV.AppID,
	}

	//下发参数缓存，对应一个设备一个记录
	duration, _ := time.ParseDuration(conf.Redis.SendRecordExpired)
	sendRecordCache:=send.SendRecordCache{}
	sendRecordCache.Init(
		conf.Redis.Server,
		conf.Redis.DB,
		duration,
		conf.Redis.Password)

	deviceSignalCache:=send.DeviceSignalCache{}
	deviceSignalCache.Init(
		conf.Redis.Server,
		conf.Redis.DeviceSignalCacheDB,
		0,
		conf.Redis.Password)

	//
	duration, _ = time.ParseDuration(conf.Redis.HeartbeatLockExpired)
	heartbeatLock:=busi.HeartbeatLock{}
	heartbeatLock.Init(
		conf.Redis.Server,
		conf.Redis.HeartbeatLockDB,
		duration,
		conf.Redis.Password)

	//实际的业务处理模块
	busiModule:=busi.Busi{
		CrvClient:&crvClinet,
		SendRecordCache:&sendRecordCache,
		HeartbeatLock:&heartbeatLock,
	}

	//busiModule.DealDeviceHeartbeat("3ec783","LSJW949UUMS997068")

	//mqttclient
	mqttClient:=mqtt.MQTTClient{
		Broker:conf.MQTT.Broker,
		User:conf.MQTT.User,
		Password:conf.MQTT.Password,
		HeartbeatTopic:conf.MQTT.HeartbeatTopic,
		DiagResponseTopic:conf.MQTT.DiagResponseTopic,
		Handler:&busiModule,
		ClientID:conf.MQTT.ClientID,
	}
	mqttClient.Init()

	//kafka consumer
	saicinterface.StartConsumer(&conf.Kafka,&crvClinet)

	//idm.InitIntegration(&conf.IDMIntegration,&crvClinet)
	duration, _ = time.ParseDuration(conf.Redis.IdmAppDataSyncLockExpired)
	idmSyncLock:=idm.IdmSyncLock{}
	idmSyncLock.Init(
		conf.Redis.Server,
		conf.Redis.IdmAppDataSyncLockDB,
		duration,
		conf.Redis.Password)

	idm.InitAppDataSyncTask(&conf.IDMIntegration,&crvClinet,&idmSyncLock)

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
    conf.Mongo.Password,
		conf.FilePath,
		&crvClinet)
	reportController.Bind(router)

	sendController:=send.SendController{
		CRVClient:&crvClinet,
		MQTTClient:&mqttClient,
		SendRecordCache:&sendRecordCache,
		DeviceSignalCache:&deviceSignalCache,
		FilePath:conf.FilePath,
		DBCUploadTopic:conf.MQTT.DBCUploadTopic,
	}
	sendController.Bind(router)

	oauthContrller:=oauth.OauthController{
		LoginUrl:conf.Oauth.Url,
	}
	oauthContrller.Bind(router)

	router.Run(conf.Service.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}