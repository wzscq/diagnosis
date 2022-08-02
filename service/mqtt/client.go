package mqtt

import (
	"log"
	"strings"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"digimatrix.com/diagnosis/busi"
)

type MQTTClient struct {
	Broker string
	User string
	Password string
	HeartbeatTopic string
	Busi *busi.Busi
}

func (mqc *MQTTClient) getClient()(*mqtt.Client){
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqc.Broker)
	opts.SetClientID("diagonsis_mqtt_subscribe_client")
	opts.SetUsername(mqc.User)
	opts.SetPassword(mqc.Password)
	opts.SetAutoReconnect(true)

	opts.SetDefaultPublishHandler(mqc.messagePublishHandler)
	opts.OnConnect = mqc.connectHandler
	opts.OnConnectionLost = mqc.connectLostHandler
	opts.OnReconnecting = mqc.reconnectingHandler

	client:=mqtt.NewClient(opts)
	if token:=client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error)
		return nil
	}
	return &client
}

func (mqc *MQTTClient) connectHandler(client mqtt.Client){
	log.Println("MQTTClient connectHandler connect status: ",client.IsConnected())
}

func (mqc *MQTTClient) connectLostHandler(client mqtt.Client, err error){
	log.Println("MQTTClient connectLostHandler connect status: ",client.IsConnected(),err)
}

func (mqc *MQTTClient) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient messagePublishHandler topic: ",msg.Topic())
}

func (mqc *MQTTClient) reconnectingHandler(Client mqtt.Client,opts *mqtt.ClientOptions){
	log.Println("MQTTClient reconnectingHandler ")
}

func (mqc *MQTTClient) onHeartbeat(Client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient onHeartbeat ",msg.Topic())
	strTopic:=msg.Topic()[len(mqc.HeartbeatTopic)-1:]
	log.Println("MQTTClient onHeartbeat strTopic ",strTopic)
	idx:=strings.Index(strTopic,":")
	deviceID:=strTopic[:idx]
	vin:=strTopic[idx+1:]
	log.Printf("MQTTClient onHeartbeat deviceID:%s,vin:%s",deviceID,vin)
	//更新心跳记录
	mqc.Busi.DealDeviceHeartbeat(deviceID,vin)
}

func (mqc *MQTTClient) Init(){
	client:=mqc.getClient()
	(*client).Subscribe(mqc.HeartbeatTopic,0,mqc.onHeartbeat)
}