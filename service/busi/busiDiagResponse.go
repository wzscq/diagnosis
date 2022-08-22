package busi

import (
	"log"
	"digimatrix.com/diagnosis/crv"
)

const (
	MODEL_SENDREC = "diag_param_sendrecord"
)

func (busi *Busi)DealDiagResponse(deviceID string){
	//获取缓存记录
	rec,err:=busi.SendRecordCache.GetSendRecord(deviceID)
	if err!=nil {
		log.Println("DealDiagResponse get cached send record error")
		return
	}

	//登录
	if busi.CrvClient.Login() ==0 {
		rec[crv.SAVE_TYPE_COLUMN]=crv.SAVE_UPDATE
		rec["status"]="1"
		//添加心跳记录到记录表
		saveReq:=&crv.CommonReq{
			ModelID:MODEL_SENDREC,
			List:&[]map[string]interface{}{
				rec,
			},
		}
		log.Println(saveReq)
		busi.CrvClient.Save(saveReq)
	}
}