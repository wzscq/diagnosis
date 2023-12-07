package busi

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/send"
)

type Busi struct {
	CrvClient *crv.CRVClient
	SendRecordCache *send.SendRecordCache
	HeartbeatLock *HeartbeatLock
}

