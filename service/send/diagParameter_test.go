package send

import (
	"fmt"
	"testing"
	"digimatrix.com/diagnosis/common"
	"digimatrix.com/diagnosis/crv"
)

func TestGetSendParameter(t *testing.T) {
	crvClient:=&crv.CRVClient{
		Server:"http://localhost:8200",
		Token:"carapi",
	}

	sp:=&SendParameter{
		CRVClient:crvClient,
		SignalList:&map[string]interface{}{},
	}

	ids:=[]string{"33"}
	dp,errorCode:=getDiagParams(ids,sp.CRVClient,sp.SignalList,"carapi")
	if errorCode!=common.ResultSuccess {
		t.Error("getSendParameter getDiagParams error")
	}

	fmt.Println(dp.Params)
}