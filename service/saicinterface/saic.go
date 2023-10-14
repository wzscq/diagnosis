package saicinterface

import (
	"digimatrix.com/diagnosis/common"
	"digimatrix.com/diagnosis/crv"
	kafka "github.com/segmentio/kafka-go"
	"context"
	"log"
	"encoding/json"
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
}

type KafkaConsumer struct {
	KafkaConf common.KafkaConf
	CRVClient *crv.CRVClient
}

const (
	MODEL_PROJECT = "diag_platform"
	MODEL_DEVICE = "vehiclemanagement"
)

func getKafkaReader(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
	})
}

func (kc *KafkaConsumer)SaveProject(pdpmPrj *PDPMProject){
	//登录
	if kc.CRVClient.Login() ==0 {
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
	}
}

func (kc *KafkaConsumer)SaveVeichle(vehicle *EVDMSVeihcle){
	
}

func (kc *KafkaConsumer)SaveDevice(device *EVDMSDevice){
	//登录
	if kc.CRVClient.Login() ==0 {
		rec:=map[string]interface{}{}
		rec[crv.SAVE_TYPE_COLUMN]=crv.SAVE_CREATE
		rec["VehicleManagementCode"]=device.VehicleNo
		rec["VIN"]=device.Vin
		rec["ProjectNum"]=device.ProjectNo
		rec["TestSpecification"]=device.Standard
		rec["DeviceNumber"]=device.DeviceCode
		rec["id"]=device.BindingDate

		//添加心跳记录到记录表
		saveReq:=&crv.CommonReq{
			ModelID:MODEL_DEVICE,
			List:&[]map[string]interface{}{
				rec,
			},
		}
		log.Println(saveReq)
		kc.CRVClient.Save(saveReq,"")
	}
}

func (kc *KafkaConsumer)ConsumePDPMProject(){
	//{"projectNo":"prjno","projectSN":"prjsn","currentPhase":"currentphase","productPlatform":"productPlatform","projectType":"projectType"}
	log.Println("start ConsumePDPMProject ... ")
	reader := getKafkaReader(kc.KafkaConf.Brokers, kc.KafkaConf.TopicPDPMProject, kc.KafkaConf.GroupID)
	defer reader.Close()
	pdpmPrj:=&PDPMProject{}
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("message at topic:%v partition:%v offset:%v  %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		err=json.Unmarshal(m.Value, &pdpmPrj)
		if err != nil {
			log.Println(err)
			continue
		}
		kc.SaveProject(pdpmPrj)
	}
}

func (kc *KafkaConsumer)ConsumeEVDMSVeihcle(){
	//{"vhicleNo":"vhicleNo","phase":"Phase","projectNo":"projectNo","vin":"Vin"}
	log.Println("start ConsumeEVDMSVeihcle ... ")
	reader := getKafkaReader(kc.KafkaConf.Brokers, kc.KafkaConf.TopicEVDMSVeihcle, kc.KafkaConf.GroupID)
	defer reader.Close()
	veihcle:=&EVDMSVeihcle{}
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("message at topic:%v partition:%v offset:%v  %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		err=json.Unmarshal(m.Value, &veihcle)
		if err != nil {
			log.Println(err)
			continue
		}
		kc.SaveVeichle(veihcle)
	}
}

func (kc *KafkaConsumer)ConsumeEVDMSDevice(){
	//{"deviceCode":"deviceCode","vehicleNo":"vehicleNo","vin":"vin","projectNo":"projectNo","standard":"standard","developPhase":"developPhase","bindingDate":"bindingDate","untieDate":"untieDate"}
	log.Println("start ConsumeEVDMSDevice ... ")
	reader := getKafkaReader(kc.KafkaConf.Brokers, kc.KafkaConf.TopicEVDMSDevice, kc.KafkaConf.GroupID)
	defer reader.Close()
	device:=&EVDMSDevice{}
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("message at topic:%v partition:%v offset:%v  %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		err=json.Unmarshal(m.Value, &device)
		if err != nil {
			log.Println(err)
			continue
		}
		kc.SaveDevice(device)
	}
}

func StartConsumer(kafkaConf common.KafkaConf,crvClient *crv.CRVClient){
	
	kafkaConsumer:=KafkaConsumer{
		KafkaConf:kafkaConf,
		CRVClient:crvClient,
	}

	go kafkaConsumer.ConsumePDPMProject()
	//go kafkaConsumer.ConsumeEVDMSVeihcle()
	//go kafkaConsumer.ConsumeEVDMSDevice() 
}