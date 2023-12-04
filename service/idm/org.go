package idm

import (
	"log"
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
)

func (appData *AppData)SyncOrgs(updateAt,appToken string){
	number:=38
	for {
		//get user info
		orgs,err:=GetAppOrgBy(appData.GetAppOrgByUrl,appToken,updateAt,number)
		if err!=nil {
			log.Println("DoSync error:",err.Error())
			return
		}
		if orgs==nil || len(orgs)==0 {
			return
		}		

		log.Println("DoSync orgs count:",len(orgs)," number:",number)
		for _,org:= range orgs {
			appData.SyncOrg(org)
		}

		number++
	}
}

func (appData *AppData)SyncOrg(org idmOrg){
	crvOrg:=appData.GetCrvOrg(org.Code)
	if crvOrg!=nil {
		log.Println("crvOrg:",crvOrg)
		appData.UpdateCRVOrg(&org,crvOrg)

	} else {
		appData.CreateCRVOrg(&org)
	}
}

func (appData *AppData)GetCrvOrg(id string)(map[string]interface{}){
	//查询数据
	commonRep:=crv.CommonReq{
		ModelID:"diag_org",
		Fields:&[]map[string]interface{}{
			{"field":"id"},
			{"field":"version"},
		},
		Filter:&map[string]interface{}{
			"id":id,
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

func (appData *AppData)UpdateCRVOrg(idmOrg *idmOrg,crvOrg map[string]interface{})(error){
	disable:="否"
	if idmOrg.IsDisabled {
		disable="是"
	}

	commonRep:=crv.CommonReq{
		ModelID:"diag_org",
		List:&[]map[string]interface{}{
			{
				"id":crvOrg["id"],
				"_save_type":"update",
				"version":crvOrg["version"],
				"disable":disable,
				"organization_id":idmOrg.OrganizationID,
				"name":idmOrg.Name,
		        "parent_id":idmOrg.ParentID,
				"org_type":idmOrg.OrgType,
				"sequence":idmOrg.Sequence,
				"zhrdxcj":idmOrg.CJ,
			},
		},
	}

	appData.CRVClient.Save(&commonRep,"")
	return nil
}

func (appData *AppData)CreateCRVOrg(idmOrg *idmOrg)(error){
	disable:="否"
	if idmOrg.IsDisabled {
		disable="是"
	}

	commonRep:=crv.CommonReq{
		ModelID:"diag_org",
		List:&[]map[string]interface{}{
			{
				"id":idmOrg.Code,
				"_save_type":"create",
				"disable":disable,
				"organization_id":idmOrg.OrganizationID,
				"name":idmOrg.Name,
		        "parent_id":idmOrg.ParentID,
				"org_type":idmOrg.OrgType,
				"sequence":idmOrg.Sequence,
				"zhrdxcj":idmOrg.CJ,
			},
		},
	}

	appData.CRVClient.Save(&commonRep,"")
	return nil
}

