package busi

import (
	"testing"
	"log"
	"os"
	"time"
	"digimatrix.com/diagnosis/mqtt"
)

func _TestHuaweiOSSClient_ListObjects(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	huaweiOSSClient := &HuaweiOSSClient{
		AccessKeyID: "HPUAZSP3LUEP3S86CUCP",
		SecretAccessKey: "HQ1MBmc7AShxJpwRiXShXUQp3LPHSiikDapRKVbX",
		EndPoint: "obs.cn-east-3.myhuaweicloud.com",
		BucketName: "obs-qa-evdms",
	}
	huaweiOSSClient.ListObjectsSDK()
}

func _TestHuaweiOSSClient_GetObjectSDK(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	huaweiOSSClient := &HuaweiOSSClient{
		AccessKeyID: "HPUAZSP3LUEP3S86CUCP",
		SecretAccessKey: "HQ1MBmc7AShxJpwRiXShXUQp3LPHSiikDapRKVbX",
		EndPoint: "obs.cn-east-3.myhuaweicloud.com",
		BucketName: "obs-qa-evdms",
		OutputPath: "D:/",
	}

	key := "app/sharefolder/4021057/rec_4021057_2023-10-29-14_52_20.dat.gz"
	huaweiOSSClient.GetObjectSDK(key)
}

func TestHuaweiOSSClient_DownloadOSSFile(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	huaweiOSSClient := &HuaweiOSSClient{
		AccessKeyID: "HPUAZSP3LUEP3S86CUCP",
		SecretAccessKey: "HQ1MBmc7AShxJpwRiXShXUQp3LPHSiikDapRKVbX",
		EndPoint: "obs.cn-east-3.myhuaweicloud.com",
		BucketName: "obs-qa-evdms",
		OutputPath: "D:/",
	}

	busiModule:=Busi{
		HuaweiOSSClient:huaweiOSSClient,
	}

	mqttClient:=mqtt.MQTTClient{
		Broker: "tcp://49.4.3.226:1983",
		User: "mosquitto",
		Password: "123456",
		DownloadOSSFileTopic: "download_oss_file",
		HeartbeatTopic:"MQB1/Status/+",
        DiagResponseTopic:"MQB1/ReciveAck/+",
        ClientID:"diagonsis_mqtt_subscribe_client_2",
		Handler:&busiModule,
	}

	mqttClient.Init()
	time.Sleep(10 * time.Second)
	mqttClient.Publish("download_oss_file", "app/sharefolder/4021057/rec_4021057_2023-10-29-14_52_20.dat.gz")

	//循环检测文件是否下载完成
	for {
		time.Sleep(1 * time.Second)
		if _, err := os.Stat("D:/rec_4021057_2023-10-29-14_52_20.dat.gz"); err == nil {
			log.Println("文件下载完成")
			break
		}
	}
}