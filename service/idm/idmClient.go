package idm

import (
	"log"
	"net/http"
	"net/url"
	"encoding/json"
	"errors"
)

func GetAppAccBy(reqUrl,token string,updateAt string)([]idmUser,error){
	searchMap:=map[string]interface{}{
		"updateAt_gt":updateAt,
	}
	searchJson,_:= json.Marshal(searchMap)
	pageMap:=map[string]interface{}{
		"size":10000,
		"number":1,
	}
	pageJson,_:=json.Marshal(pageMap)

	data := url.Values{}
  data.Set("token", token)
  data.Set("_search", string(searchJson))
	data.Set("_page", string(pageJson))

	rsp, err := http.PostForm(reqUrl, data)
	if err != nil {
		log.Println("GetAppAccBy error:",err.Error())
		return nil,err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 { 
		log.Println("GetAppAccBy rsp.StatusCode:",rsp.StatusCode)
		err := errors.New("GetAppAccBy rsp.StatusCode")
		return nil,err
	}
	
	decoder := json.NewDecoder(rsp.Body)
	var appAccByRsp getAppAccByRsp
	err = decoder.Decode(&appAccByRsp)
	if err != nil {
		log.Println("GetAppAccBy error:",err.Error())
		return nil,err
	}

	return appAccByRsp.Data,nil
}

func GetAppToken(reqUrl,systemCode,integrationKey,clientID string)(string,error){
	data := url.Values{}
  data.Set("systemCode", systemCode)
  data.Set("integrationKey", integrationKey)
	data.Set("force", "true")
	data.Set("clientId", clientID)

	rsp, err := http.PostForm(reqUrl, data)
	if err != nil {
		log.Println("GetAppToken error:",err.Error())
		return "",err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 { 
		log.Println("GetAppToken rsp.StatusCode:",rsp.StatusCode)
		err := errors.New("GetAppToken rsp.StatusCode")
		return "",err
	}
	
	decoder := json.NewDecoder(rsp.Body)
	var appTokenRsp getAppTokenRsp
	err = decoder.Decode(&appTokenRsp)
	if err != nil {
		log.Println("GetAppToken error:",err.Error())
		return "",err
	}

	log.Println(appTokenRsp)

	return appTokenRsp.Data,nil

}