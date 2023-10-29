package oauth

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

type OauthController struct {
	LoginUrl string
}

func (controller *OauthController)redirectLogin(c *gin.Context){
	state:=GetBatchID()
	//重定向web到给定的回调地址
	url:=controller.LoginUrl+state
	log.Println("end OAuthController login Redirect to",url)
	c.Redirect(http.StatusMovedPermanently, url)
}	

func (controller *OauthController) Bind(router *gin.Engine) {
	log.Println("Bind OauthController")
	router.GET("/redirectLogin",controller.redirectLogin)
}