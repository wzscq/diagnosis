package mqtt

import (
	"log"
	"strings"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"digimatrix.com/diagnosis/common"
	"github.com/google/uuid"
	"time"
)

const (
	MSG_TYPE_DIAG="Diag"
	MSG_TYPE_EVENT="Event"
	MSG_TYPE_SIGNAL="SignalFilter"
)

type eventHandler interface {
	DealDeviceHeartbeat(deviceID,vin string)
	DealDiagResponse(deviceID string)
	DealEventResponse(deviceID string)
	DealSignalResponse(deviceID string)
	DealDownloadOSSFile(key string)
}

type MQTTClient struct {
	Broker string
	User string
	Password string
	HeartbeatTopic string
	DiagResponseTopic string
	DownloadOSSFileTopic string
	ClientID string
	Handler eventHandler
	Client mqtt.Client
}

func (mqc *MQTTClient) getClient()(mqtt.Client){
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqc.Broker)
	clientID:=uuid.New().String()
	opts.SetClientID(clientID)
	opts.SetUsername(mqc.User)
	opts.SetPassword(mqc.Password)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(time.Second*10)

	opts.SetDefaultPublishHandler(mqc.messagePublishHandler)
	opts.OnConnect = mqc.connectHandler
	opts.OnConnectionLost = mqc.connectLostHandler
	opts.OnReconnecting = mqc.reconnectingHandler

	client:=mqtt.NewClient(opts)
	if token:=client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return nil
	}
	return client
}

func (mqc *MQTTClient) connectHandler(client mqtt.Client){
	log.Println("MQTTClient connectHandler connect status: ",client.IsConnected())
	if client.IsConnected() {
		log.Println("MQTTClient Subscribe HeartbeatTopic: ",mqc.HeartbeatTopic)
		client.Subscribe(mqc.HeartbeatTopic,0,mqc.onHeartbeat)
		log.Println("MQTTClient Subscribe DiagResponseTopic: ",mqc.DiagResponseTopic)
		client.Subscribe(mqc.DiagResponseTopic,0,mqc.onDiagResponse)
		log.Println("MQTTClient Subscribe DownloadOSSFileTopic: ",mqc.DownloadOSSFileTopic)
		client.Subscribe(mqc.DownloadOSSFileTopic,0,mqc.onDownloadOSSFile)
	}
}

func (mqc *MQTTClient) connectLostHandler(client mqtt.Client, err error){
	log.Println("MQTTClient connectLostHandler connect status: ",client.IsConnected(),err)
}

func (mqc *MQTTClient) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient messagePublishHandler topic: ",msg.Topic())
}

func (mqc *MQTTClient) reconnectingHandler(client mqtt.Client,opts *mqtt.ClientOptions){
	log.Println("MQTTClient reconnectingHandler ")
}

func (mqc *MQTTClient)onDownloadOSSFile(Client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient onDownloadOSSFile ",msg.Topic())
	key:=string(msg.Payload())
	log.Println("MQTTClient onDownloadOSSFile key: ",key)
	mqc.Handler.DealDownloadOSSFile(key)
}

func (mqc *MQTTClient)onHeartbeat(Client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient onHeartbeat ",msg.Topic())
	strTopic:=msg.Topic()[len(mqc.HeartbeatTopic)-1:]
	log.Println("MQTTClient onHeartbeat strTopic ",strTopic)
	idx:=strings.Index(strTopic,":")
	deviceID:=strTopic[:idx]
	vin:=strTopic[idx+1:]
	log.Printf("MQTTClient onHeartbeat deviceID:%s,vin:%s",deviceID,vin)
	//更新心跳记录
	mqc.Handler.DealDeviceHeartbeat(deviceID,vin)
}

func (mqc *MQTTClient)onDiagResponse(Client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient onDiagResponse ",msg.Topic())
	strTopic:=msg.Topic()[len(mqc.DiagResponseTopic)-1:]
	log.Println("MQTTClient onDiagResponse strTopic ",strTopic)
	idx:=strings.Index(strTopic,":")
	deviceID:=strTopic[:idx]
	strType:=strTopic[idx+1:]
	log.Printf("MQTTClient onDiagResponse deviceID:%s,type:%s",deviceID,strType)
	//更新下发状态
	if strType==MSG_TYPE_DIAG {
		mqc.Handler.DealDiagResponse(deviceID)
	}

	if strType==MSG_TYPE_EVENT {
		mqc.Handler.DealEventResponse(deviceID)
	}

	if strType==MSG_TYPE_SIGNAL {
		mqc.Handler.DealSignalResponse(deviceID)
	}
}

func (mqc *MQTTClient)Publish(topic,content string)(int){
	if mqc.Client == nil {
		return common.ResultMqttClientError
	}
	log.Println("MQTTClient Publish topic:"+topic+" content:"+content)
	token:=mqc.Client.Publish(topic,0,false,content)
	token.Wait()
	return common.ResultSuccess
}

func (mqc *MQTTClient) Init(){
	mqc.Client=mqc.getClient()
	if mqc.Client == nil {
		return
	}
}