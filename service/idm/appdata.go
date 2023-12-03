package idm

import (
	"log"
	"time"
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
)

type getAppTokenRsp struct {
	Success bool `json:"success"`
	Data string `json:"data"`
	ErrorMessage string `json:"errorMessage"`
}

type getAppAccByRsp struct {
	Success bool `json:"success"`
	Data []idmUser `json:"data"`
}

type AppData struct {
	GetAppTokenUrl string
	GetAppAccBy string
	SystemCode string
	IntegrationKey string
	ClientID string
	Duration string
	CRVClient *crv.CRVClient
	RoleMap map[string]string
	DefaultRole string
	InitUpdateAt string
	UpdateTime string
}

func InitAppDataSyncTask(integrationConf *common.IntegrationConf,crvClient *crv.CRVClient){
	appData := &AppData{
		GetAppTokenUrl:integrationConf.GetAppTokenUrl,
		GetAppAccBy:integrationConf.GetAppAccBy,
		SystemCode: integrationConf.SystemCode,
		IntegrationKey: integrationConf.IntegrationKey,
		ClientID:integrationConf.ClientID,
		Duration: integrationConf.Duration,
		RoleMap: integrationConf.RoleMap,
		DefaultRole: integrationConf.DefaultRole,
		CRVClient:crvClient,
		InitUpdateAt:integrationConf.InitUpdateAt,
		UpdateTime:integrationConf.UpdateTime,
	}

	//i.TestDealUserAccount()
	
	go appData.start()
}

func (appData *AppData) start() {
	//第一次同步
	if appData.InitUpdateAt!="" {
		log.Println("AppData init sync start ...")
		appData.DoSync(appData.InitUpdateAt)
		log.Println("AppData init sync end")
	}

	//获取最新的更新时间
	updateAt:=time.Now().Format("2006-01-02 15:04:05")

	//第一次同步后等待一段时间，然后开始定时同步
	//计算当前时间到下一个更新时间的时间间隔
	duration, _ := time.ParseDuration(appData.Duration)
	if appData.UpdateTime!="" {
		now := time.Now()
		updateTime,_ := time.Parse("15:04:05",appData.UpdateTime)
		updateTime = time.Date(now.Year(),now.Month(),now.Day(),updateTime.Hour(),updateTime.Minute(),updateTime.Second(),0,time.Local)
		for {
			if updateTime.Before(now) {
				updateTime = updateTime.Add(duration)
			} else {
				break
			}
		}
		duration := updateTime.Sub(now)
		log.Println("AppData sync duration:",duration)
		time.Sleep(duration)
	}

	for {
		currentUpdateAt := updateAt
		//当前时间减去一分钟，防止同步过程中有新的数据
		updateAt = time.Now().Add(-time.Minute*1).Format("2006-01-02 15:04:05")
		log.Println("AppData sync start with updateAt:",currentUpdateAt," ...")
		appData.DoSync(currentUpdateAt)
		log.Println("AppData sync end")
		time.Sleep(duration)
	}
}

func (appData *AppData) DoSync(updateAt string) {
	//get app token
	appToken,err := GetAppToken(appData.GetAppTokenUrl,appData.SystemCode,appData.IntegrationKey,appData.ClientID)
	if err!=nil {
		log.Println("DoSync error:",err.Error())
		return
	}

	//get user info
	users,err:=GetAppAccBy(appData.GetAppAccBy,appToken,updateAt)
	if err!=nil {
		log.Println("DoSync error:",err.Error())
		return
	}

	for _,user:= range users {
		appData.UpdateUser(user)
	}

	log.Println("DoSync users:",users)
}

func (appData *AppData)UpdateUser(user idmUser){
	//获取用户角色
	
}

