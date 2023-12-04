package idm

import (
	"log"
	"net/http"
	"net/url"
	"encoding/json"
	"errors"
)

type idmUser struct {
	UserName string `json:"username"`
	FullName string `json:"fullname"`
	Email string `json:"email"`
	RoleList []string `json:"roleList"`
	Mobile string `json:"mobile"`
	UserDepartShort string `json:"userDepartShort"`
	IsLocked bool `json:"isLocked"`
	IsDisabled bool `json:"isDisabled"`
	AID string `json:"_AID"`
	ZHRSSGW string `json:"ZHRSSGW"`
	OrganizationID string `json:"organizationId"`
}

type idmOrg struct {
	OrganizationID string `json:"organizationId"`
	Name string `json:"name"`
	ParentID string `json:"parentId"`
	OrgType string `json:"orgType"`
	Code string `json:"code"`
	Sequence int `json:"sequence"`
	IsDisabled bool `json:"isDisabled"`
	CJ string `json:"ZHRDXCJ"`
}

type getAppTokenRsp struct {
	Success bool `json:"success"`
	Data string `json:"data"`
	ErrorMessage string `json:"errorMessage"`
}

type getAppAccByRsp struct {
	Success bool `json:"success"`
	Data []idmUser `json:"data"`
}

type getAppOrgByRsp struct {
	Success bool `json:"success"`
	Data []idmOrg `json:"data"`
}

func GetAppAccBy(reqUrl,token,updateAt string,number int)([]idmUser,error){
	searchMap:=map[string]interface{}{
		"updateAt_gt":updateAt,
	}
	searchJson,_:= json.Marshal(searchMap)
	pageMap:=map[string]interface{}{
		"size":100,
		"number":number,
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

func GetAppOrgBy(reqUrl,token,updateAt string,number int)([]idmOrg,error){
	searchMap:=map[string]interface{}{
		"updateAt_gt":updateAt,
	}
	searchJson,_:= json.Marshal(searchMap)
	pageMap:=map[string]interface{}{
		"size":100,
		"number":number,
	}
	pageJson,_:=json.Marshal(pageMap)

	data := url.Values{}
  	data.Set("token", token)
  	data.Set("_search", string(searchJson))
	data.Set("_page", string(pageJson))

	rsp, err := http.PostForm(reqUrl, data)
	if err != nil {
		log.Println("GetAppOrgBy error:",err.Error())
		return nil,err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 { 
		log.Println("GetAppOrgBy rsp.StatusCode:",rsp.StatusCode)
		err := errors.New("GetAppOrgBy rsp.StatusCode")
		return nil,err
	}
	
	decoder := json.NewDecoder(rsp.Body)
	var appOrgByRsp getAppOrgByRsp
	err = decoder.Decode(&appOrgByRsp)
	if err != nil {
		log.Println("GetAppOrgBy error:",err.Error())
		return nil,err
	}

	return appOrgByRsp.Data,nil
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