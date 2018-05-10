package core

import (
	"os"
	"testing"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/configuration"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/hte"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp/securenet"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types/event"
)

func TestNewSDKFactory(t *testing.T) {

	sdkFactory, err := NewSDKFactory()
	if err != nil {
		t.Error("Error on geting new SDK factory.")
		t.FailNow()
	}
	if sdkFactory == nil {
		t.Error("Failed to get new SDK factory.")
		t.FailNow()
	}
}

func createUUIDFile(t *testing.T, fileContent string) {

	var file, err = os.OpenFile(UUIDFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		t.Error("Cannot open new file: " + UUIDFilePath + ", " + err.Error())
		t.FailNow()
	}
	defer file.Close()

	_, err = file.WriteString(fileContent)

	if err != nil {
		t.Error("Cannot write to file: " + UUIDFilePath + ", " + err.Error())
		t.FailNow()
	}
	err = file.Sync()
	if err != nil {
		t.Error("Cannot sync file, " + err.Error())
		t.FailNow()
	}
}

func removeUUIDFile(t *testing.T) {
	err := os.Remove(UUIDFilePath)
	if err != nil {
		t.Error("Cannot remove file, " + err.Error())
	}
}

func testGetDeviceWithoutExistingUUIDFile(t *testing.T) {

	sdkFactory, err := NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}

	device, err := sdkFactory.GetDevice("test", "test", "test", &configuration.WPWithin{})
	defer removeUUIDFile(t)

	// will not error, it will create new uuid file
	if err != nil {
		t.Error("Should not error, should create new uuid file.")
		t.FailNow()
	}
	if device == nil {
		t.Error("Device should not be nil for this case.")
		t.FailNow()
	}

}

func testGetDeviceWithExistingUUIDFile(t *testing.T) {

	sdkFactory, err := NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}

	createUUIDFile(t, "405aae3a-9ed1-11e7-abc4-cec278b6b50a")
	defer removeUUIDFile(t)

	device, err := sdkFactory.GetDevice("test", "test", "test", &configuration.WPWithin{})
	if err != nil {
		t.Error("Error on geting device.")
		t.FailNow()
	}
	if device == nil {
		t.Error("Device is nil.")
		t.FailNow()
	}
}

func TestGetDevice(t *testing.T) {

	testGetDeviceWithoutExistingUUIDFile(t)

	testGetDeviceWithExistingUUIDFile(t)

}

func TestGetPSPMerchant(t *testing.T) {

	var err error
	var sdkFactory SDKFactory
	sdkFactory, err = NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}

	var psp psp.PSP
	psp, err = sdkFactory.GetPSPMerchant(nil)
	if err == nil {
		t.Error("Should fail for nil argument.")
		t.FailNow()
	}
	if psp != nil {
		t.Error("psp should be nil for error case scenario.")
		t.FailNow()
	}

	var mymap map[string]string
	mymap = make(map[string]string)
	mymap["dummy"] = "some dummy string"
	psp, err = sdkFactory.GetPSPMerchant(mymap)
	if err == nil {
		t.Error("Should fail for dummy map.")
		t.FailNow()
	}
	if psp != nil {
		t.Error("psp should be nil for dummy map.")
		t.FailNow()
	}

	mymap = make(map[string]string)
	mymap["psp_name"] = onlineworldpay.PSPName
	psp, err = sdkFactory.GetPSPMerchant(mymap)
	if err != nil {
		t.Error("Shouldn't fail for onlineworldpay.PSPName.")
		t.FailNow()
	}
	if psp == nil {
		t.Error("psp shouldn't be nil for onlineworldpay.PSPName.")
		t.FailNow()
	}

	mymap["psp_name"] = securenet.PSPName
	mymap[securenet.CfgDeveloperID] = "abc"
	psp, err = sdkFactory.GetPSPMerchant(mymap)
	if err == nil {
		t.Error("Should fail abc is not an integer.")
		t.FailNow()
	}
	if psp != nil {
		t.Error("psp should be nil for dummy map.")
		t.FailNow()
	}
	mymap[securenet.CfgDeveloperID] = "123"
	psp, err = sdkFactory.GetPSPMerchant(mymap)
	if err != nil {
		t.Error("Shouldn't fail.")
		t.FailNow()
	}
	if psp == nil {
		t.Error("psp shouldn't be nil.")
		t.FailNow()
	}
}

func TestGetPSPClient(t *testing.T) {
	var err error
	var sdkFactory SDKFactory
	sdkFactory, err = NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}

	var psp0 psp.PSP
	psp0, err = sdkFactory.GetPSPClient(nil)
	if err == nil {
		t.Error("Should fail for nil argument.")
		t.FailNow()
	}
	if psp0 != nil {
		t.Error("psp should be nil for error case scenario.")
		t.FailNow()
	}

	var pspCfg map[string]string
	pspCfg = make(map[string]string)

	pspCfg[psp.CfgPSPName] = "dummy"

	psp0, err = sdkFactory.GetPSPClient(pspCfg)
	if err == nil {
		t.Error("Should fail for dummy map.")
		t.FailNow()
	}
	if psp0 != nil {
		t.Error("psp should be nil for dummy map.")
		t.FailNow()
	}

	pspCfg[psp.CfgPSPName] = "worldpayonlinepayments"
	psp0, err = sdkFactory.GetPSPClient(pspCfg)
	if err != nil {
		t.Error("Shouldn't fail for worldpayonlinepayments.")
		t.FailNow()
	}
	if psp0 == nil {
		t.Error("psp shouldn't be nil for dummy map.")
		t.FailNow()
	}
	pspCfg[psp.CfgPSPName] = "securenet"
	psp0, err = sdkFactory.GetPSPClient(pspCfg)
	if err == nil {
		t.Error("Should fail for securenet.")
		t.FailNow()
	}
	if psp0 != nil {
		t.Error("psp should be nil for dummy map.")
		t.FailNow()
	}
	pspCfg[securenet.CfgDeveloperID] = "123"
	psp0, err = sdkFactory.GetPSPClient(pspCfg)
	if err != nil {
		t.Error("Shouldn't fail for securenet if pspCfg[securenet.CfgDeveloperID] is set.")
		t.FailNow()
	}
	if psp0 == nil {
		t.Error("psp shouldn't be nil for securenet if pspCfg[securenet.CfgDeveloperID] is set.")
		t.FailNow()
	}
}

func TestGetSvcBroadcaster(t *testing.T) {
	var err error
	var sdkFactory SDKFactory
	sdkFactory, err = NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}
	broadcaster0, err := sdkFactory.GetSvcBroadcaster("127.0.0.1")
	if err != nil {
		t.Error("Error on geting broadcaster.")
		t.FailNow()
	}
	if broadcaster0 == nil {
		t.Error("GetSvcBroadcaster returns nil broadcaster.")
		t.FailNow()
	}
}

func TestGetSvcScanner(t *testing.T) {
	var err error
	var sdkFactory SDKFactory
	sdkFactory, err = NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}
	scanner0, err := sdkFactory.GetSvcScanner()
	if err != nil {
		t.Error("Error on geting scanner.")
		t.FailNow()
	}
	if scanner0 == nil {
		t.Error("GetSvcScanner returns nil.")
		t.FailNow()
	}
}

func TestGetHTE(t *testing.T) {
	sdkFactory, err := NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}

	var psp0 psp.PSP
	var om0 hte.OrderManager
	hteService0, err := sdkFactory.GetHTE(nil, psp0, "127.0.0.1", "scheme", nil, om0, nil)

	if err != nil {
		t.Error("Error on geting HTE service.")
		t.FailNow()
	}
	if hteService0 == nil {
		t.Error("GetHTE returns nil.")
		t.FailNow()
	}
}

func TestGetOrderManager(t *testing.T) {
	sdkFactory, err := NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}
	var omng0 hte.OrderManager
	omng0, err = sdkFactory.GetOrderManager()
	if err != nil {
		t.Error("Error on geting order manager.")
		t.FailNow()
	}
	if omng0 == nil {
		t.Error("GetOrderManager returns nil.")
		t.FailNow()
	}
}

func TestGetHTEClient(t *testing.T) {
	sdkFactory, err := NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}
	hteClient0, err := sdkFactory.GetHTEClient()
	if hteClient0 != nil || err != nil {
		t.Error("GetHTEClient returns nil, because is not implemented yet.")
		t.FailNow()
	}
}

func TestGetHTEClientHTTP(t *testing.T) {
	sdkFactory, err := NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}
	hteClientHTTP0, err := sdkFactory.GetHTEClientHTTP()
	if err != nil {
		t.Error("Error on return from GetHTEClientHTTP.")
		t.FailNow()
	}
	if hteClientHTTP0 == nil {
		t.Error("GetHTEClientHTTP() returns nil client.")
		t.FailNow()
	}
}

func TestGetHTEServiceHandler(t *testing.T) {
	sdkFactory, err := NewSDKFactory()
	if err != nil {
		t.Error("Error on geting SDK factory.")
		t.FailNow()
	}
	var device0 *types.Device
	var psp0 psp.PSP
	var credential0 *hte.Credential
	var orderManager0 hte.OrderManager
	var eventHandler0 event.Handler
	hteServHandler0 := sdkFactory.GetHTEServiceHandler(device0, psp0, credential0, orderManager0, eventHandler0)
	if hteServHandler0 == nil {
		t.Error("GetHTEServiceHandler() returns nil.")
		t.FailNow()
	}
}
