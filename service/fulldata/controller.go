package fulldata

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"io"
	"github.com/gin-gonic/gin"
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
)

type Controller struct {
	CRVClient *crv.CRVClient
}

//Bind bind the controller function to url
func (controller *Controller) Bind(router *gin.Engine) {
	log.Println("Bind controller")
	router.POST("/fulldata/download",controller.download)
	router.POST("/fulldata/sendConf",controller.sendConf)
}

func (controller *Controller) sendConf(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end redirect with error")
		return
	}	
	
	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
  	}	

	if rep.SelectedRowKeys==nil || len(*rep.SelectedRowKeys)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//获取参数信息
	dataID:=(*rep.SelectedRowKeys)[0]
	log.Println(dataID)
	conf,err:=GetConf(header.Token,dataID,controller.CRVClient)
	if err!=nil {
		log.Println(err)
		params:=map[string]interface{}{
			"errorMsg": err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetFullDataConfError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	log.Println(conf.Name)

	//获取车辆信息
	if rep.List == nil || len(*rep.List) == 0 {
		params:=map[string]interface{}{
			"errorMsg": "缺少车辆信息",
		}
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	firstItem := (*rep.List)[0]
	cars, ok := firstItem["cars"]
	if !ok || cars == nil {
		params := map[string]interface{}{
			"errorMsg": "缺少车辆信息",
		}
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	carsList,ok:=cars.(map[string]interface{})["list"]
	if !ok || carsList == nil {
		params := map[string]interface{}{
			"errorMsg": "缺少车辆列表信息",
		}
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	
	carsArray, ok := carsList.([]interface{})
	if !ok || len(carsArray) == 0 {
		params := map[string]interface{}{
			"errorMsg": "车辆列表格式错误或为空",
		}
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	for _, carItem := range carsArray {
		car, ok := carItem.(map[string]interface{})
		if !ok {
			params := map[string]interface{}{
				"errorMsg": "车辆信息格式错误",
			}
			rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}
		
		// 处理车辆信息
		carID, ok := car["id"].(string)
		if !ok {
			params := map[string]interface{}{
				"errorMsg": "车辆ID格式错误",
			}
			rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}
		
		// 可以根据需要继续处理车辆的其他信息
		log.Printf("处理车辆信息: %s", carID)
		carInfo,err:=GetCar(header.Token,carID,controller.CRVClient)
		if err!=nil {
			params:=map[string]interface{}{
				"errorMsg": err.Error(),
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultGetCarInfoError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}

		log.Println(carInfo)

		//创建下发记录
		err=CreateSendRec(header.Token,carInfo,conf,controller.CRVClient)
		if err!=nil {
			params:=map[string]interface{}{
				"errorMsg": err.Error(),
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultCreateFullDataSendRecError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}
	}

	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (controller *Controller) download(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end redirect with error")
		return
	}	
	
	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
  	}	

	if rep.SelectedRowKeys==nil || len(*rep.SelectedRowKeys)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	dataID:=(*rep.SelectedRowKeys)[0]
	//获取报告路径
	filePath,commonErr:=GetFullDataFileName(dataID,header.Token,controller.CRVClient)
	if commonErr!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(commonErr,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	log.Println(filePath)
	
	file, err := os.Open(filePath)
	if err!=nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultOpenFileError,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	defer file.Close()
	_,fileName:=filepath.Split(filePath)
	c.Header("Content-Type", "application/octet-stream")
  	c.Header("Content-Disposition", "attachment; filename="+fileName)
  	c.Header("Content-Transfer-Encoding", "binary")
	io.Copy(c.Writer,file)
}