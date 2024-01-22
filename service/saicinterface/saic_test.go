package saicinterface

import (
	"testing"
	"digimatrix.com/diagnosis/crv"
	"encoding/json"
)

func TestKafkaConsumer(t *testing.T) {
	CRVClient := &crv.CRVClient{
		Server: "http://127.0.0.1:8200",
		Token:"carapi",
	}
	kc := &KafkaConsumer{
		CRVClient: CRVClient,
	}

	device:=EVDMSDevice{
		DeviceCode:"test01",
		VehicleNo:"vehicleNo",
		Vin:"vin001",
		ProjectNo:"prj001",
		Standard:"test01",
		DevelopPhase:"test01",
		BindingDate:"2024-01-20 00:00:00",
		UntieDate:"2024-03-20 00:00:00",
		VehicleConfiger:"1",
	}

	deviceMsg:=EVDMSDeviceMsg{
		DataType:"1",
		Number:1,
		Detail:[]EVDMSDevice{
			device,
		},
	}

	deviceStr,_:=json.Marshal(deviceMsg)
	kc.SaveDevice(string(deviceStr))
}