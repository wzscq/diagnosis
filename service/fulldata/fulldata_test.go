package fulldata

import (
	"testing"
	"digimatrix.com/diagnosis/crv"
	"log"
	"digimatrix.com/diagnosis/common"
)

func TestGetFullDataFileName(t *testing.T) {
	CRVClient := &crv.CRVClient{
		Server: "http://127.0.0.1:8200",
		Token:"carapiv2",
	}
	
	filename,errcode:=GetFullDataFileName("32","carapiv2",CRVClient)
	if errcode!=common.ResultSuccess {
		t.Errorf("GetFullDataFileName failed")
	}
	log.Println(filename)	
}