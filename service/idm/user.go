package idm

import (
	"log"
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
)

func (appData *AppData)SyncUsers(updateAt,appToken string){
	//get user info
	number:=1
	for {
		users,err:=GetAppAccBy(appData.GetAppAccByUrl,appToken,updateAt,number)
		if err!=nil {
			log.Println("DoSync error:",err.Error())
			return
		}

		if users==nil || len(users)==0 {
			return
		}

		for _,user:= range users {
			appData.SyncUser(&user)
		}
		
		number++
	}
}

func (appData *AppData)SyncUser(user *idmUser){
	log.Println("SyncUser:",*user)
	//获取用户部门
	user.OrganizationID=appData.getUserOrg(user.OrganizationID)
	//获取用户角色

	//查询本地用户信息
	crvUser:=appData.GetCRVUserInfo(user.UserName)
	if crvUser!=nil {
		appData.UpdateCRVUser(user,crvUser)
	} else {
		appData.CreateCRVUser(user)
	}
}

func (appData *AppData)GetCRVUserInfo(userID string)(map[string]interface{}){
	//查询数据
	commonRep:=crv.CommonReq{
		ModelID:"core_user",
		Fields:&userInfoFields,
		Filter:&map[string]interface{}{
			"id":userID,
		},
	}

	req,commonErr:=appData.CRVClient.Query(&commonRep,"")
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

func (appData *AppData)CreateCRVUser(idmUser *idmUser)(error){
	//查询数据
  	/*roleList:=[]map[string]interface{}{}
	for _,role:=range idmUser.RoleList {
		roleList=append(roleList,map[string]interface{}{
			"id":role,
			"_save_type":"create",
		})
	}*/

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
				"department":idmUser.OrganizationID,
				"job_number":idmUser.AID,
				"dimission":dimission,
				"disable":disable,
				"password":"a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
				"_save_type":"create",
				/*"roles":map[string]interface{}{
					"list":roleList,
					"fieldType":"many2many",
					"modelID":"core_role",
				},*/
			},
		},
	}
	appData.CRVClient.Save(&commonRep,"")
	return nil
}

func (appData *AppData)UpdateCRVUser(idmUser *idmUser,crvUser map[string]interface{})(error){
	/*roleList:=[]map[string]interface{}{}
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
	}*/

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
				"department":idmUser.OrganizationID,
				"job_number":idmUser.AID,
				"dimission":dimission,
				"disable":disable,
				/*"roles":map[string]interface{}{
					"list":roleList,
					"fieldType":"many2many",
					"modelID":"core_role",
				},*/
			},
		},
	}

	appData.CRVClient.Save(&commonRep,"")
	return nil
}

func (appData *AppData)getUserOrg(orgID string)(string){
	//查询数据
	commonRep:=crv.CommonReq{
		ModelID:"diag_org",
		Fields:&[]map[string]interface{}{
			map[string]interface{}{"field":"id"},
			map[string]interface{}{"field":"name"},
			map[string]interface{}{"field":"zhrdxcj"},
			map[string]interface{}{
				"field":"parent_id",
				"fieldType":"many2one",
				"relatedModelID":"diag_org",
				"fields":&[]map[string]interface{}{
					map[string]interface{}{"field":"id"},
					map[string]interface{}{"field":"name"},
					map[string]interface{}{"field":"zhrdxcj"},
					map[string]interface{}{
						"field":"parent_id",
						"fieldType":"many2one",
						"relatedModelID":"diag_org",
						"fields":&[]map[string]interface{}{
							map[string]interface{}{"field":"id"},
							map[string]interface{}{"field":"name"},
							map[string]interface{}{"field":"zhrdxcj"},
							map[string]interface{}{
								"field":"parent_id",
								"fieldType":"many2one",
								"relatedModelID":"diag_org",
								"fields":&[]map[string]interface{}{
									map[string]interface{}{"field":"id"},
									map[string]interface{}{"field":"name"},
									map[string]interface{}{"field":"zhrdxcj"},
									map[string]interface{}{
										"field":"parent_id",
										"fieldType":"many2one",
										"relatedModelID":"diag_org",
										"fields":&[]map[string]interface{}{
											map[string]interface{}{"field":"id"},
											map[string]interface{}{"field":"name"},
											map[string]interface{}{"field":"zhrdxcj"},
										},
									},		
								},
							},
						},
					},		
				},
			},
		},
		Filter:&map[string]interface{}{
			"id":orgID,
		},
	}

	req,commonErr:=appData.CRVClient.Query(&commonRep,"")
	if commonErr!=common.ResultSuccess {
		return orgID
	}

	if req.Error == true {
		log.Println("GetProjectData error:",req.ErrorCode,req.Message)
		return orgID
	}

	if req.Result["list"]!=nil && len(req.Result["list"].([]interface{}))>0 {
		log.Println("getUserOrg:",req.Result["list"].([]interface{})[0].(map[string]interface{}))
		return GetOrgName(req.Result["list"].([]interface{})[0].(map[string]interface{}))	
	}
	
	log.Println("getUserOrg error:orgID=",orgID)
	
	return orgID
}

func GetOrgName(org map[string]interface{})(string){
	//获取机构层级
	orgCj:=org["zhrdxcj"].(string)
	orgName:=org["name"].(string)
	//如果当前层级小于等于40，直接返回
	if orgCj<="40" {
		return orgName		
	}

	//如果当前层级大于40，获取上级机构
	parentOrg:=org["parent_id"].(map[string]interface{})
	if parentOrg["list"]!=nil && len(parentOrg["list"].([]interface{}))>0 {
		return GetOrgName(parentOrg["list"].([]interface{})[0].(map[string]interface{}))
	}

	return orgName
}