package saicinterface

import (
	"digimatrix.com/diagnosis/common"
	"digimatrix.com/diagnosis/crv"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"encoding/json"
	"time"
)

type PDPMProject struct {
	ProjectNo string `json:"projectNo"`
	ProjectSN string `json:"projectSN"`
	CurrentPhase string `json:"currentPhase"`
	ProductPlatform string `json:"productPlatform"`
	ProjectType string `json:"projectType"`
}


type EVDMSVeihcle struct {
	VhicleNo string `json:"vhicleNo"`
	Phase string `json:"phase"`
	ProjectNo string `json:"projectNo"`
	Vin string `json:"vin"`
}

type EVDMSDevice struct {
	DeviceCode string `json:"deviceCode"`
	VehicleNo string `json:"vehicleNo"`
	Vin string `json:"vin"`
	ProjectNo string `json:"projectNo"`
	Standard string `json:"standard"`
	DevelopPhase string `json:"developPhase"`
	BindingDate string `json:"bindingDate"`
	UntieDate string `json:"untieDate"`
	VehicleConfiger string `json:"vehicleConfiger"`
}

type EVDMSDeviceMsg struct {
	DataType string `json:"dataType"`
	Detail []EVDMSDevice `json:"detail"`
	Number int `json:"number"`
}

type KafkaConsumer struct {
	KafkaConf *common.KafkaConfig
	CRVClient *crv.CRVClient
}

const (
	MODEL_PROJECT = "diag_platform"
	MODEL_DEVICE = "vehiclemanagement"
)

func (kc *KafkaConsumer)SaveProject(projectString string){
	pdpmPrj:=&PDPMProject{}
	err:=json.Unmarshal([]byte(projectString),pdpmPrj)
	if err != nil {
		log.Println(err)
		return
	}

	//登录
	//if kc.CRVClient.Login() ==0 {
		rec:=map[string]interface{}{}
		rec[crv.SAVE_TYPE_COLUMN]=crv.SAVE_CREATE
		rec["id"]=pdpmPrj.ProjectNo
		rec["name"]=pdpmPrj.ProjectSN
		//添加心跳记录到记录表
		saveReq:=&crv.CommonReq{
			ModelID:MODEL_PROJECT,
			List:&[]map[string]interface{}{
				rec,
			},
		}
		log.Println(saveReq)
		kc.CRVClient.Save(saveReq,"")
	//}
}

func (kc *KafkaConsumer)SaveVeichle(vehicle *EVDMSVeihcle){
	
}

func (kc *KafkaConsumer)SaveDevice(deviceString string){
	device:=&EVDMSDeviceMsg{}
	err:=json.Unmarshal([]byte(deviceString),device)
	if err != nil {
		log.Println(err)
		return
	}

	if len(device.Detail)==0 {
		log.Println("no device detail")
		return
	}

	//登录
	//if kc.CRVClient.Login() ==0 {
		reclst:=[]map[string]interface{}{}
		for _,deviceItem:=range device.Detail {
			rec:=map[string]interface{}{}
			rec[crv.SAVE_TYPE_COLUMN]=crv.SAVE_CREATE
			rec["VehicleManagementCode"]=deviceItem.VehicleNo
			rec["VIN"]=deviceItem.Vin
			rec["ProjectNum"]=deviceItem.ProjectNo
			rec["TestSpecification"]=deviceItem.Standard
			rec["DeviceNumber"]=deviceItem.DeviceCode
			rec["id"]=deviceItem.DeviceCode+"_"+deviceItem.BindingDate
			rec["BindingDate"]=deviceItem.BindingDate
			rec["UntieDate"]=deviceItem.UntieDate
			rec["developPhase"]=deviceItem.DevelopPhase
			rec["vehicleConfiger"]=deviceItem.VehicleConfiger
			reclst=append(reclst,rec)
		}
		//添加心跳记录到记录表
		saveReq:=&crv.CommonReq{
			ModelID:MODEL_DEVICE,
			List:&reclst,
		}
		log.Println(saveReq)
		kc.CRVClient.Save(saveReq,"")
		
	//}
}

func (kc *KafkaConsumer)doInitConsumer() *kafka.Consumer {
	log.Print("init kafka consumer, it may take a few seconds to init the connection\n")
	//common arguments
	var kafkaconf = &kafka.ConfigMap{
			"api.version.request": "true",
			"auto.offset.reset": "earliest",
			"heartbeat.interval.ms": 3000,
			"session.timeout.ms": 30000,
			"max.poll.interval.ms": 120000,
			"fetch.max.bytes": 1024000,
			"max.partition.fetch.bytes": 256000,
	}
	kafkaconf.SetKey("bootstrap.servers", kc.KafkaConf.BootstrapServers);
	kafkaconf.SetKey("group.id", kc.KafkaConf.GroupId)

	switch kc.KafkaConf.SecurityProtocol {
	case "PLAINTEXT" :
			kafkaconf.SetKey("security.protocol", "plaintext");
	case "SASL_SSL":
			kafkaconf.SetKey("security.protocol", "sasl_ssl");
			kafkaconf.SetKey("ssl.ca.location", "./conf/ca-cert.pem");
			kafkaconf.SetKey("sasl.username", kc.KafkaConf.SaslUsername);
			kafkaconf.SetKey("sasl.password", kc.KafkaConf.SaslPassword);
			kafkaconf.SetKey("sasl.mechanism", kc.KafkaConf.SaslMechanism)
	case "SASL_PLAINTEXT":
			kafkaconf.SetKey("security.protocol", "sasl_plaintext");
			kafkaconf.SetKey("sasl.username", kc.KafkaConf.SaslUsername);
			kafkaconf.SetKey("sasl.password", kc.KafkaConf.SaslPassword);
			kafkaconf.SetKey("sasl.mechanism", kc.KafkaConf.SaslMechanism)

	default:
		  log.Println("init kafka consumer error:","unknown protocol")
			//panic(kafka.NewError(kafka.ErrUnknownProtocol, "unknown protocol", true))
			return nil
	}

	consumer, err := kafka.NewConsumer(kafkaconf)
	if err != nil {
		log.Println("init kafka consumer error:",err)
		return nil
	}
	log.Print("init kafka consumer success\n")
	return consumer;
}

func (kc *KafkaConsumer)Start(){
	log.Println("KafkaConsumer start ... ")	
	consumer := kc.doInitConsumer()
	if consumer == nil {
		return
	}

	log.Println("KafkaConsumer start SubscribeTopics",kc.KafkaConf.TopicPDPMProject,kc.KafkaConf.TopicEVDMSDevice)	
	consumer.SubscribeTopics([]string{kc.KafkaConf.TopicPDPMProject,kc.KafkaConf.TopicEVDMSDevice}, nil)
	for {
			log.Println("KafkaConsumer start ReadMessage ...")
			msg, err := consumer.ReadMessage(time.Second * 10)
			if err == nil {
				log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				switch *msg.TopicPartition.Topic {
				case kc.KafkaConf.TopicPDPMProject:
					kc.SaveProject(string(msg.Value))
				case kc.KafkaConf.TopicEVDMSDevice:
					kc.SaveDevice(string(msg.Value))
				}
			} else {
				// The client will automatically try to recover from all errors.
				log.Printf("KafkaConsumer error: %v (%v)\n", err, msg)
			}
	}
	consumer.Close()
	log.Println("KafkaConsumer Closed")
}

func StartConsumer(kafkaConf *common.KafkaConfig,crvClient *crv.CRVClient){
	
	kafkaConsumer:=KafkaConsumer{
		KafkaConf:kafkaConf,
		CRVClient:crvClient,
	}

	go kafkaConsumer.Start()
	//go kafkaConsumer.ConnectToKafa()
	//go kafkaConsumer.ConsumePDPMProject()
	//go kafkaConsumer.ConsumeEVDMSVeihcle()
	//go kafkaConsumer.ConsumeEVDMSDevice() 
}

