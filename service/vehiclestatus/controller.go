package vehiclestatus

import (
	"log"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
)

type VehicleStatusController struct {
	uri    string
	isOpen bool
	client *mongo.Client
	dbName string
	Username  string
	Password  string
}

func CreateController(uri, dbName,userName,password string) *VehicleStatusController {
	return &VehicleStatusController{
		uri: uri, 
		isOpen: false, 
		client: nil, 
		dbName: dbName,
		Username:userName,
		Password:password,
	}
}

//Open establish a connection to mongoDB server.
func (controller *VehicleStatusController) Init() error {
	log.Println("VehicleStatusController open ...")
	auth:=options.Credential{
		AuthSource:"admin",
		Username:controller.Username,
		Password:controller.Password,
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(controller.uri).SetAuth(auth))
	if err != nil {
		log.Println("VehicleStatusController create mongo client error ", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Println("VehicleStatusController connect to mongodb server error ", err)
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("VehicleStatusController can not connect to mongodb server,error ", err)
		return err
	}

	log.Println("VehicleStatusController connect to mongodb ok")
	controller.client = client
	controller.isOpen = true
	return nil
}

//Close disconnect from mongoDB server.
func (controller *VehicleStatusController) Release() error {
	if controller.isOpen {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		controller.client.Disconnect(ctx)
		controller.isOpen = false
		controller.client = nil
		log.Println("VehicleStatusController disconnect from mongodb")
	}
	return nil
}

func (controller *VehicleStatusController) getVehicelStatus (c *gin.Context){
	statusid := c.Param("statusid")
	log.Println("getVehicelStatus statusid:",statusid)
	params:=strings.Split(statusid, "_")
	var jsonDocuments []map[string]interface{}
	if len(params)<2 {
		log.Println("getVehicelStatus statusid is wrong")
		c.JSON(http.StatusOK, jsonDocuments)
		return
	}
	timeStamp,err:= strconv.ParseInt(params[1],10,32)
	if err!=nil {
		log.Println("getVehicelStatus can not convert timeStamp to int32")
		c.JSON(http.StatusOK, jsonDocuments)
		return
	}
	collectionName:=params[0]
	//convert timeStampe to time
	stime:=time.Unix(timeStamp, 0).Format("2006-01-02 15:04:05")
	
	log.Println("getVehicelStatus stime:",stime," collectionName:",collectionName)

	filter:=bson.D{
		{"SamplingTime",
			bson.D{{"$gte", stime},{"$lte", stime}},
		},
	}

	controller.Init()
	defer controller.Release()

	if controller.isOpen != true {
		c.JSON(http.StatusOK, jsonDocuments)
		return
	}

	collection := controller.client.Database(controller.dbName).Collection(collectionName)
	ctxFind, cancelFind := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancelFind()
	cur, errFind := collection.Find(ctxFind, filter,nil)
	if errFind != nil {
		log.Println("getVehicelStatus find documents error")
		log.Println(errFind)
		c.JSON(http.StatusOK, jsonDocuments)
		return
	}
	defer cur.Close(ctxFind)
	var bsonDocument bson.D
	var temporaryBytes []byte
	for cur.Next(ctxFind) {
		cur.Decode(&bsonDocument)
		temporaryBytes, _ = bson.MarshalExtJSON(bsonDocument, true, true)
		var jsonDocument map[string]interface{}
		json.Unmarshal(temporaryBytes, &jsonDocument)
		jsonDocuments = append(jsonDocuments, jsonDocument)
	}
	c.JSON(http.StatusOK, jsonDocuments)
}

//Bind bind the controller function to url
func (controller *VehicleStatusController) Bind(router *gin.Engine) {
	log.Println("Bind VehicleStatusController")
	router.GET("/vehicleStatus/:statusid", controller.getVehicelStatus)
}