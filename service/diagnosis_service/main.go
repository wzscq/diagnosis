package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"digimatrix.com/diagnosis/dashboard"
	"digimatrix.com/diagnosis/common"
	"digimatrix.com/diagnosis/report"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowAllOrigins:true,
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
        AllowCredentials: true,
    }))

	conf:=common.InitConfig()
	
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
        conf.Mongo.Password)
	reportController.Bind(router)

	router.Run(conf.Service.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}