package dashboard

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type Controller struct {
	Repository Repository
}

type DashBoradFilter struct {
	Project string `json:"project"`
	Year string `json:"year"`
	Type string `json:"type"`
	Spec string `json:"spec"`
}

func GetWhereString(filter DashBoradFilter)(string){
	str:="where 1=1 "
	if len(filter.Year)>0 {
		str=str+"and time like '"+filter.Year+"%' "
	}

	if len(filter.Project)>0 {
		str=str+"and project_num like '"+filter.Project+"' "
	}

	if len(filter.Type)>0 {
		str=str+"and type like '%"+filter.Type+"%' "
	}

	if len(filter.Spec)>0 {
		str=str+"and specifications like '"+filter.Spec+"' "
	}

	return str
}

func (controller *Controller) getDashboard(c *gin.Context) {
	var filter DashBoradFilter
	if err := c.BindJSON(&filter); err != nil {
		log.Println(err)
		return
  }	

	log.Printf("getDashboard with filters: year:%s, project:%s,type: %s,spec:%s\n",filter.Year,filter.Project,filter.Type,filter.Spec)
	whereStr:=GetWhereString(filter)
	carCount,_:=controller.Repository.getCarCount(whereStr)
	projectcarCount,_:=controller.Repository.getCarCountByProject(whereStr)
	faultCountByType,_:=controller.Repository.getFaultCountByType(whereStr)
	faultCountByStatus,_:=controller.Repository.getFaultCountByStatus(whereStr)
	faultList,_:=controller.Repository.getFaultList(whereStr)
	
	res := map[string]interface{}{
		"carCount":     carCount,
		"projectcarCount":projectcarCount,
		"faultCountByType":faultCountByType,
		"faultCountByStatus":faultCountByStatus,
		"faultList":faultList,
	}
	c.JSON(http.StatusOK, res)
}

func (controller *Controller) getDashboardYears(c *gin.Context) {
	var filter DashBoradFilter
	if err := c.BindJSON(&filter); err != nil {
		log.Println(err)
		return
  }	

	log.Printf("getDashboard with filters: year:%s, project:%s,type: %s,spec:%s\n",filter.Year,filter.Project,filter.Type,filter.Spec)
	whereStr:=GetWhereString(filter)

	years,_:=controller.Repository.getDashboardYears(whereStr)
	c.JSON(http.StatusOK, years)
}

func (controller *Controller) getDashboardProjects(c *gin.Context) {
	var filter DashBoradFilter
	if err := c.BindJSON(&filter); err != nil {
		log.Println(err)
		return
  }	

	log.Printf("getDashboard with filters: year:%s, project:%s,type: %s,spec:%s\n",filter.Year,filter.Project,filter.Type,filter.Spec)
	whereStr:=GetWhereString(filter)

	projects,_:=controller.Repository.getDashboardProjects(whereStr)
	c.JSON(http.StatusOK, projects)
}

func (controller *Controller) getDashboardTypes(c *gin.Context) {
	var filter DashBoradFilter
	if err := c.BindJSON(&filter); err != nil {
		log.Println(err)
		return
  }	

	log.Printf("getDashboard with filters: year:%s, project:%s,type: %s,spec:%s\n",filter.Year,filter.Project,filter.Type,filter.Spec)
	whereStr:=GetWhereString(filter)

	projects,_:=controller.Repository.getDashboardTypes(whereStr)
	c.JSON(http.StatusOK, projects)
}

func (controller *Controller) getDashboardSpecs(c *gin.Context) {
	var filter DashBoradFilter
	if err := c.BindJSON(&filter); err != nil {
		log.Println(err)
		return
  }	

	log.Printf("getDashboard with filters: year:%s, project:%s,type: %s,spec:%s\n",filter.Year,filter.Project,filter.Type,filter.Spec)
	whereStr:=GetWhereString(filter)

	projects,_:=controller.Repository.getDashboardSpecs(whereStr)
	c.JSON(http.StatusOK, projects)
}

func (controller *Controller) Bind(router *gin.Engine) {
	log.Println("Bind MysqlController")
	router.POST("/dashboard", controller.getDashboard)
	router.POST("/dashboard/years", controller.getDashboardYears)
	router.POST("/dashboard/projects", controller.getDashboardProjects)
	router.POST("/dashboard/types", controller.getDashboardTypes)
	router.POST("/dashboard/specs", controller.getDashboardSpecs)
}