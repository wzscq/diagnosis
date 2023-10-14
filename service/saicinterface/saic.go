package saicinterface

import (
	"digimatrix.com/diagnosis/common"
	kafka "github.com/segmentio/kafka-go"
	"context"
	"log"
)

type KafkaConsumer struct {
	KafkaConf common.KafkaConf
}

func getKafkaReader(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
	})
}

func (kc *KafkaConsumer)ConsumePDPMProject(){
	log.Println("start ConsumePDPMProject ... ")
	reader := getKafkaReader(kc.KafkaConf.Brokers, kc.KafkaConf.TopicPDPMProject, kc.KafkaConf.GroupID)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("message at topic:%v partition:%v offset:%v  %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}

func StartConsumer(kafkaConf common.KafkaConf){
	kafkaConsumer:=KafkaConsumer{
		KafkaConf:kafkaConf,
	}
	go kafkaConsumer.ConsumePDPMProject()
}