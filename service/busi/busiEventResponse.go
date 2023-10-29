package busi

import (
	"log"
	"digimatrix.com/diagnosis/crv"
)

const (
	MODEL_EVNET_SENDREC = "diag_event_sendrecord"
)

func (busi *Busi)DealEventResponse(deviceID string){
	//获取缓存记录
	rec,err:=busi.SendRecordCache.GetSendRecord("event_"+deviceID)
	if err!=nil {
		log.Println("DealEventResponse get cached send record error")
		return
	}

	//登录
	//if busi.CrvClient.Login() ==0 {
		rec[crv.SAVE_TYPE_COLUMN]=crv.SAVE_UPDATE
		rec["status"]="1"
		//添加心跳记录到记录表
		saveReq:=&crv.CommonReq{
			ModelID:MODEL_EVNET_SENDREC,
			List:&[]map[string]interface{}{
				rec,
			},
		}
		log.Println(saveReq)
		busi.CrvClient.Save(saveReq,"")
	//}
}