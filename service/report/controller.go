package report

import (
	"log"
	"context"
	"strconv"
	"encoding/json"
	"time"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	uri    string
	isOpen bool
	client *mongo.Client
	dbName string
	Username  string
	Password  string
}

//CreateController create a controller
func CreateController(uri, dbName,userName,password string) *Controller {
	return &Controller{
		uri: uri, 
		isOpen: false, 
		client: nil, 
		dbName: dbName,
		Username:userName,
		Password:password,
	}
}

//Open establish a connection to mongoDB server.
func (controller *Controller) Init() error {
	log.Println("controller open ...")
	auth:=options.Credential{
		AuthSource:"admin",
		Username:controller.Username,
		Password:controller.Password,
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(controller.uri).SetAuth(auth))
	if err != nil {
		log.Println("create mongo client error ", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Println("connect to mongodb server error ", err)
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("can not connect to mongodb server,error ", err)
		return err
	}

	log.Println("connect to mongodb ok")
	controller.client = client
	controller.isOpen = true
	return nil
}

//Close disconnect from mongoDB server.
func (controller *Controller) Release() error {
	if controller.isOpen {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		controller.client.Disconnect(ctx)
		controller.isOpen = false
		controller.client = nil
		log.Println("disconnect from mongodb")
	}
	return nil
}

func (controller *Controller) getReports(c *gin.Context){
	log.Println("getReports start")
	var jsonDocuments []map[string]interface{}
	var countDocument int64
	var err = error(nil)
	countDocument=0
	collectionName:=c.Query("collection")
	dtc:=c.Query("dtc")
	if(collectionName!=""&&dtc!=""){
		controller.Init()
		defer controller.Release()

		if controller.isOpen == true {
			startDate:=c.DefaultQuery("startDate","")
			endDate:=c.DefaultQuery("endDate","")
			pageStr:=c.DefaultQuery ("page","0")
			countStr:=c.DefaultQuery("count","10")
			var page int64
			page, err = strconv.ParseInt(pageStr,10,64)
			if(err!=nil){
				page=0
			}
			var count int64
			count, err = strconv.ParseInt(countStr,10,64)
			if(err!=nil){
				count=10
			}
			opts := options.Find().SetLimit(count).SetSkip(page*count)

			filter:=bson.D{
				{"SamplingTime",
					bson.D{{"$gte", startDate},{"$lte", endDate}},
				},
				{
				    "ReportID",dtc,
				},
			}

			log.Println("getReports count documents ...")
			collection := controller.client.Database(controller.dbName).Collection(collectionName)
			ctxCount, cancelCount := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancelCount()
			resultCount,errCount :=collection.CountDocuments(ctxCount,filter)
			if(errCount!=nil){
				log.Println("getReports count documents error")
				log.Println(errCount)
				err = errCount
			} else {
				countDocument=resultCount
				log.Println("getReports find documents ...")
				ctxFind, cancelFind := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancelFind()
				cur, errFind := collection.Find(ctxFind, filter,opts)
				defer cur.Close(ctxFind)
				if errFind != nil {
					log.Println("getReports find documents error")
					log.Println(errFind)
					err = errFind
				} else {
					var bsonDocument bson.D
					var temporaryBytes []byte
					for cur.Next(ctxFind) {
						err = cur.Decode(&bsonDocument)
						temporaryBytes, err = bson.MarshalExtJSON(bsonDocument, true, true)
						var jsonDocument map[string]interface{}
						err = json.Unmarshal(temporaryBytes, &jsonDocument)
						jsonDocuments = append(jsonDocuments, jsonDocument)
					}
					if err = cur.Err(); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
	log.Println("getReports over")
	res := map[string]interface{}{
		"total":     countDocument,
		"data":      jsonDocuments,
	}
	c.JSON(http.StatusOK, res)
}

func (controller *Controller) downloadReport (c *gin.Context){
	var repo ReportContent
	if err := c.BindJSON(&repo); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusOK, "")
		return
    }
	report:=repo.getReport()
	c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Disposition", "attachment; filename="+repo.FileName+".xlsx")
    c.Header("Content-Transfer-Encoding", "binary")
	report.Write(c.Writer)
}

//Bind bind the controller function to url
func (controller *Controller) Bind(router *gin.Engine) {
	log.Println("Bind controller")
	router.POST("/reports", controller.getReports)
	router.POST("/downloadReport",controller.downloadReport)
}