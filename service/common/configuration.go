package common

import (
	"log"
	"os"
	"encoding/json"
)

type mysqlConf struct {
	Server string `json:"server"`
	Password string `json:"password"`
	User string `json:"user"`
	DBName string `json:"dbName"`
}

type mongoConf struct {
	Server string `json:"server"`
	Password string `json:"password"`
	User string `json:"user"`
	DBName string `json:"dbName"`
}

type serviceConf struct {
	Port string `json:"port"`
}

type mqttConf struct {
	Broker string `json:"broker"`
	User string `json:"user"`
	Password string `json:"password"`
	HeartbeatTopic string `json:"heartbeatTopic"`
	DiagResponseTopic string `json:"diagResponseTopic"`
	DBCUploadTopic string `json:"dbcUploadTopic"`
	ClientID string `json:"clientID"`
	DownloadOSSFileTopic string `json:"downloadOSSFileTopic"`
}

type crvConf struct {
	Server string `json:"server"`
  //User string `json:"user"`
  //Password string `json:"password"`
  AppID string `json:"appID"`
	Token string `json:"token"`
}

type redisConf struct {
	Server string `json:"server"`
  SendRecordExpired string `json:"sendRecordExpired"`
  DB int `json:"db"`
	DeviceSignalCacheDB int 
	Password string `json:"password"`
	HeartbeatLockDB int  `json:"heartbeatLockDB"`
	HeartbeatLockExpired string `json:"heartbeatLockExpired"`
	IdmAppDataSyncLockDB int  `json:"idmAppDataSyncLockDB"`
	IdmAppDataSyncLockExpired string `json:"idmAppDataSyncLockExpired"`
}

type IntegrationConf struct {
	Url string `json:"url"`
	GetAppTokenUrl string `json:"getAppTokenUrl"`
	GetAppAccByUrl string `json:"getAppAccByUrl"`
	GetAppOrgByUrl string `json:"getAppOrgByUrl"`
	ClientID string `json:"clientID"`
	SystemCode string `json:"systemCode"`
	IntegrationKey string `json:"integrationKey"`
	Duration string `json:"duration"`
	RoleMap map[string]string `json:"roleMap"`
	DefaultRole string `json:"defaultRole"`
	InitUpdateAt string `json:"initUpdateAt"`
	UpdateTime string `json:"updateTime"`
}

type KafkaConfig struct {
	Brokers []string `json:"brokers"`
	TopicPDPMProject string `json:"topic.pdpm.project"`
	TopicEVDMSVeihcle string `json:"topic.evdms.veihcle"`
	TopicEVDMSDevice string `json:"topic.evdms.device"`
	GroupId    string `json:"group.id"`
	BootstrapServers    string `json:"bootstrap.servers"`
	SecurityProtocol string `json:"security.protocol"`
	SaslMechanism string `json:"sasl.mechanism"`
	SaslUsername string `json:"sasl.username"`
	SaslPassword string `json:"sasl.password"`
}

type OauthConf struct {
	Url string `json:"url"`
}

type FullData struct {
	TempPath string `json:"tempPath"`
	URLPrefix string `json:"urlPrefix"`
	TempFileExpired string `json:"tempFileExpired"`
}

type HuaweiOSSConf struct {
	AccessKeyID string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	EndPoint string `json:"endPoint"`
	BucketName string `json:"bucketName"`
	OutputPath string `json:"outputPath"`
}

type Config struct {
	Mysql  mysqlConf  `json:"mysql"`
	Service serviceConf `json:"service"`
	Mongo mongoConf `json:"mongo"`
	MQTT mqttConf `json:"mqtt"`
	CRV crvConf `json:"crv"`
	Redis redisConf `json:"redis"`
	FilePath string `json:"filePath"`
	Kafka KafkaConfig `json:"kafka"`
	Oauth OauthConf `json:"oauth"`
	IDMIntegration IntegrationConf `json:"idmIntegration"`
	VehicleStatusMongo mongoConf `json:"vehicleStatusMongo"`
	FullData FullData `json:fullData`
	HuaweiOSS HuaweiOSSConf `json:"huaweiOSS"`
}

var gConfig Config

func InitConfig()(*Config){
	log.Println("init configuation start ...")
	//获取用户账号
	//获取用户角色信息
	//根据角色过滤出功能列表
	fileName := "conf/conf.json"
	filePtr, err := os.Open(fileName)
	if err != nil {
        log.Fatal("Open file failed [Err:%s]", err.Error())
    }
    defer filePtr.Close()

	// 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&gConfig)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
	}
	log.Println("init configuation end")
	return &gConfig
}

func GetConfig()(*Config){
	return &gConfig
}