package idm

import (
	"time"
	"net/http"
	"encoding/json"
	"net/url"
	"errors"
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"log"
)

const (
	OBJ_TARGET_ACCOUNT = "TARGET_ACCOUNT"
	EF_ON_CREATE = "CREATED"
	EF_ON_UPDATE = "UPDATED"
)

type loginRsp struct {
	Success bool `json:"success"`
	TokenId string `json:"tokenId"`
}

type pullRsp struct {
	Success bool `json:"success"`
	TaskId string `json:"taskId"`
	ObjectType string `json:"objectType"`
	EffectOn string `json:"effectOn"`
	GUID string `json:"guid"`
	Data map[string]interface{} `json:"data"`
}

type Integration struct {
	Url string
	SystemCode string
	IntegrationKey string
	Duration string
	CRVClient *crv.CRVClient
	RoleMap map[string]string
	DefaultRole string
}

var userInfoFields=[]map[string]interface{}{
	{"field":"id"},
	{"field":"version"},
	{
		"field":"roles",
		"fieldType":"many2many",
		"relatedModelID":"core_role",
		"fields":[]map[string]interface{}{
			{"field":"id"},
			{"field":"version"},
		},
	},
}

func (i *Integration)GetIdmUser(pullRsp *pullRsp)(*idmUser){
	if (pullRsp.Data==nil) {
		return nil
	}

	idmUser:=&idmUser{}

	userName:=pullRsp.Data["username"]
	log.Println("username:",userName)
	if userName!=nil {
		idmUser.UserName=userName.(string)
	}
	fullname:=pullRsp.Data["fullname"]
	log.Println("username:",userName)
	if fullname!=nil {
		idmUser.FullName=fullname.(string)
	}
	email:=pullRsp.Data["email"]
	if email!=nil {
		idmUser.Email=email.(string)
	}
	mobile:=pullRsp.Data["mobile"]
	if mobile!=nil {
		idmUser.Mobile=mobile.(string)
	}
	userDepartShort:=pullRsp.Data["userDepartShort"]
	if userDepartShort!=nil {
		idmUser.UserDepartShort=userDepartShort.(string)
	}
	isLocked:=pullRsp.Data["isLocked"]
	if isLocked!=nil {
		idmUser.IsLocked=isLocked.(bool)
	}
	isDisabled:=pullRsp.Data["isDisabled"]
	if isDisabled!=nil {
		idmUser.IsDisabled=isDisabled.(bool)
	}
	ZHRSSGW:=pullRsp.Data["ZHRSSGW"]
	if ZHRSSGW!=nil {
		idmUser.ZHRSSGW=ZHRSSGW.(string)
	}

	log.Println("idm user:",userName)
	//获取角色信息
	roleList:=[]string{}
	roles:=pullRsp.Data["roleList"]
	if roles!=nil {
		roleRecords:=roles.([]interface{})
		for _,roleRecord:=range roleRecords {
			role:=roleRecord.(string)
			log.Println("idm role:",role)
			localRole,ok:=i.RoleMap[role]
			if ok {
				roleList=append(roleList,localRole)
			}
		}
	}

	if len(roleList)==0 {
		roleList=append(roleList,i.DefaultRole)
	}

	idmUser.RoleList=roleList

	return idmUser
}

func (i *Integration)GetCRVUserInfo(userID string)(map[string]interface{}){
	//查询数据
	commonRep:=crv.CommonReq{
		ModelID:"core_user",
		Fields:&userInfoFields,
		Filter:&map[string]interface{}{
			"id":userID,
		},
	}

	req,commonErr:=i.CRVClient.Query(&commonRep,"")
	if commonErr!=common.ResultSuccess {
		return nil
	}

	if req.Error == true {
		log.Println("GetProjectData error:",req.ErrorCode,req.Message)
		return nil
	}

	if req.Result["list"]!=nil && len(req.Result["list"].([]interface{}))>0 {
		return req.Result["list"].([]interface{})[0].(map[string]interface{})
	}
	return nil
}

func (i *Integration)CreateCRVUser(idmUser *idmUser)(error){
	//查询数据
  roleList:=[]map[string]interface{}{}
	for _,role:=range idmUser.RoleList {
		roleList=append(roleList,map[string]interface{}{
			"id":role,
			"_save_type":"create",
		})
	}

	dimission:="否"
	if idmUser.IsLocked {
		dimission="是"
	}

	disable:="否"
	if idmUser.IsDisabled {
		disable="是"
	}

	commonRep:=crv.CommonReq{
		ModelID:"core_user",
		List:&[]map[string]interface{}{
			{
				"id":idmUser.UserName,
				"user_name_en":idmUser.FullName,
				"user_name_zh":idmUser.FullName,
				"email":idmUser.Email,
				"department":idmUser.UserDepartShort,
				"job_number":idmUser.ZHRSSGW,
				"dimission":dimission,
				"disable":disable,
				"password":"a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
				"_save_type":"create",
				"roles":map[string]interface{}{
					"list":roleList,
					"fieldType":"many2many",
					"modelID":"core_role",
				},
			},
		},
	}
	i.CRVClient.Save(&commonRep,"")
	return nil
}

func (i *Integration)UpdateCRVUser(idmUser *idmUser,crvUser map[string]interface{})(error){
	roleList:=[]map[string]interface{}{}
	//删除没有的角色
	crvRoles:=crvUser["roles"].(map[string]interface{})["list"].([]interface{})
	for _,crvRole:=range crvRoles {
		crvRoleID:=crvRole.(map[string]interface{})["id"].(string)
		hasRole:=false
		for _,userRole:=range idmUser.RoleList {
			if crvRoleID==userRole {
				hasRole=true
				break
			}
		}
		if hasRole==false {
			roleList=append(roleList,map[string]interface{}{
				"id":crvRoleID,
				"version":crvRole.(map[string]interface{})["version"],
				"_save_type":"delete",
			})
		}
	}

	//添加新的角色
	for _,userRole:=range idmUser.RoleList {
		hasRole:=false
		for _,crvRole:=range crvRoles {
			if userRole==crvRole.(map[string]interface{})["id"].(string) {
				hasRole=true
				break
			}
		}
		if hasRole==false {
			roleList=append(roleList,map[string]interface{}{
				"id":userRole,
				"_save_type":"create",
			})
		}
	}

	dimission:="否"
	if idmUser.IsLocked {
		dimission="是"
	}

	disable:="否"
	if idmUser.IsDisabled {
		disable="是"
	}

	commonRep:=crv.CommonReq{
		ModelID:"core_user",
		List:&[]map[string]interface{}{
			{
				"id":crvUser["id"],
				"_save_type":"update",
				"version":crvUser["version"],
				"user_name_en":idmUser.FullName,
				"user_name_zh":idmUser.FullName,
				"email":idmUser.Email,
				"department":idmUser.UserDepartShort,
				"job_number":idmUser.ZHRSSGW,
				"dimission":dimission,
				"disable":disable,
				"roles":map[string]interface{}{
					"list":roleList,
					"fieldType":"many2many",
					"modelID":"core_role",
				},
			},
		},
	}

	i.CRVClient.Save(&commonRep,"")
	return nil
}

func (i *Integration)DealUserAccount(pullRsp *pullRsp)(error){
	//获取IDM用户信息
	idmUser:=i.GetIdmUser(pullRsp)
	if idmUser==nil {
		log.Println("GetIdmUser error idmUser is nil")
		return nil;
	}

	log.Println("idmUser:",idmUser)

	//查询本地用户信息
	crvUser:=i.GetCRVUserInfo(idmUser.UserName)
	if crvUser!=nil {
		log.Println("crvUser:",crvUser)
		return i.UpdateCRVUser(idmUser,crvUser)
	}
		
	return i.CreateCRVUser(idmUser)
}

func (i *Integration) DealTaskRsp(pullRsp *pullRsp)(error){
	switch pullRsp.ObjectType {
	case OBJ_TARGET_ACCOUNT:
		return i.DealUserAccount(pullRsp)
	}

	return nil
}

func (i *Integration) DoIntegration(){
	tokenId,err:=IntegrationLogin(i.Url,i.SystemCode,i.IntegrationKey)
	if err!=nil {
		log.Println("IntegrationLogin error:",err.Error())
		return
	}

	for {
		pullRsp,err:=IntegrationPullTask(i.Url,i.SystemCode,tokenId)
		if err!=nil {
			log.Println("IntegrationPullTask error:",err.Error())
			return
		}

		if !pullRsp.Success {
			log.Println("IntegrationPullTask error: pullRsp.Success is false")
			return
		}

		err=i.DealTaskRsp(pullRsp)

		if err!=nil {
			log.Println("DealTaskRsp error:",err.Error())
			return
		}
		
		err=IntegrationPullFinish(i.Url,i.SystemCode,tokenId,pullRsp.TaskId,pullRsp.GUID)
		if err!=nil {
			log.Println("IntegrationPullFinish error:",err.Error())
			return
		}
	}
}

func (i *Integration) start() {
	duration, _ := time.ParseDuration(i.Duration)
	for {
		log.Println("DoIntegration start ...")
		i.DoIntegration()
		log.Println("DoIntegration end")
		time.Sleep(duration)
	}
}

func (i *Integration)TestDealUserAccount(){
	pullRsp:=&pullRsp{
		Success:true,
		TaskId:"1",
		ObjectType:"target_account",
		EffectOn:"create",
		GUID:"1",
		Data:map[string]interface{}{
			"username":"test",
			"fullname":"test",
			"email":"test@126.com",
			"userDepartShort":"test",
			"ZHRSSGW":"test",
			"isLocked":false,
			"isDisabled":false,
			"roleList":[]interface{}{"role1"},
		},
	}

	i.DealUserAccount(pullRsp)
}

func InitIntegration(integrationConf *common.IntegrationConf,crvClient *crv.CRVClient){
	i := &Integration{
		Url: integrationConf.Url,
		SystemCode: integrationConf.SystemCode,
		IntegrationKey: integrationConf.IntegrationKey,
		Duration: integrationConf.Duration,
		RoleMap: integrationConf.RoleMap,
		DefaultRole: integrationConf.DefaultRole,
		CRVClient:crvClient,
	}

	//i.TestDealUserAccount()
	
	go i.start()
}

func IntegrationLogin(idmUrl,systemCode,integrationKey string)(string,error){
	requestData := map[string]interface{}{
		"systemCode":systemCode,
		"integrationKey":integrationKey,
		"force":true,
		"timestamp":time.Now().Unix(),
	}
	requestDataJson,_:= json.Marshal(requestData)
	data := url.Values{}
  data.Set("method", "login")
  data.Set("request", string(requestDataJson))

	log.Println("IntegrationLogin requestDataJson:",string(requestDataJson),",idmUrl:",idmUrl)

	rsp, err := http.PostForm(idmUrl, data)
	if err != nil {
		log.Println("IntegrationLogin error:",err.Error())
		return "",err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 { 
		log.Println("IntegrationLogin rsp.StatusCode:",rsp.StatusCode)
		err := errors.New("IntegrationLogin rsp.StatusCode")
		return "",err
	}
	
	decoder := json.NewDecoder(rsp.Body)
	var loginRsp loginRsp
	err = decoder.Decode(&loginRsp)
	if err != nil {
		log.Println("IntegrationLogin error:",err.Error())
		return "",err
	}

	return loginRsp.TokenId,nil
}

func IntegrationPullTask(idmUrl,systemCode,tokenId string)(*pullRsp,error){
	requestData := map[string]interface{}{
		"systemCode":systemCode,
		"tokenId":tokenId,
		"timestamp":time.Now().Unix(),
	}
	requestDataJson,_:= json.Marshal(requestData)

	data := url.Values{}
  data.Set("method", "pullTask")
  data.Set("request", string(requestDataJson))

	rsp, err := http.PostForm(idmUrl, data)
	if err != nil {
		log.Println("IntegrationLogin error:",err.Error())
		return nil,err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 { 
		log.Println("IntegrationLogin rsp.StatusCode:",rsp.StatusCode)
		err := errors.New("IntegrationLogin rsp.StatusCode")
		return nil,err
	}
	
	decoder := json.NewDecoder(rsp.Body)
	var pullRsp pullRsp
	err = decoder.Decode(&pullRsp)
	if err != nil {
		log.Println("IntegrationLogin error:",err.Error())
		return nil,err
	}
	log.Println("IntegrationLogin end with pullRsp:",pullRsp)
	return &pullRsp,nil
}

func IntegrationPullFinish(idmUrl,systemCode,tokenId,taskId,guid string)(error){
	requestData := map[string]interface{}{
		"systemCode":systemCode,
		"tokenId":tokenId,
		"taskId":taskId,
		"success":true,
		"message":"",
		"guid":guid,
		"timestamp":time.Now().Unix(),
	}
	requestDataJson,_:= json.Marshal(requestData)

	data := url.Values{}
  data.Set("method", "pullFinish")
  data.Set("request", string(requestDataJson))

	rsp, err := http.PostForm(idmUrl, data)
	if err != nil {
		log.Println("IntegrationLogin error:",err.Error())
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 { 
		log.Println("IntegrationLogin rsp.StatusCode:",rsp.StatusCode)
		err := errors.New("IntegrationLogin rsp.StatusCode")
		return err
	}
	
	return nil
}