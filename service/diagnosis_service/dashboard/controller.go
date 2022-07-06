package dashboard

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type Controller struct {
	Repository Repository
}

func (controller *Controller) getDashboard(c *gin.Context) {
	carCount,_:=controller.Repository.getCarCount()
	projectcarCount,_:=controller.Repository.getCarCountByProject()
	faultCountByType,_:=controller.Repository.getFaultCountByType()
	faultCountByStatus,_:=controller.Repository.getFaultCountByStatus()
	faultList,_:=controller.Repository.getFaultList()
	
	res := map[string]interface{}{
		"carCount":     carCount,
		"projectcarCount":projectcarCount,
		"faultCountByType":faultCountByType,
		"faultCountByStatus":faultCountByStatus,
		"faultList":faultList,
	}
	c.JSON(http.StatusOK, res)
}

func (controller *Controller) Bind(router *gin.Engine) {
	log.Println("Bind MysqlController")
	router.GET("/dashboard", controller.getDashboard)
}