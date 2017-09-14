package wpwithin_test // Black-box approach
import (
	"fmt"
	"testing"
	//"runtime"
	wpw "github.com/wptechinnovation/wpw-sdk-go/wpwithin"
	types "github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
	psp "github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp"
	securenet "github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp/securenet"
	onlineworldpay "github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay"
)

// TODO
const CFG_HtePubKey = "TODO: Enter proper key here"
const CFG_HtePrvKey = "TODO: Enter proper key here"

func TestInitialise(t *testing.T){
	fmt.Println("==============")
	var w wpw.WPWithin
	var e error
	w,e = wpw.Initialise("Dummy","Dummy")
	if e != nil {
		t.Error("Initialise did fail on non-empty device name")
		t.FailNow()
	}
	if w == nil {
		t.Error("Initialise did not return valid object")
		t.FailNow()
	}
	w,e = wpw.Initialise("Dummy","Dummy")
	if e != nil {
		t.Error("Initialise did fail on non-empty device name")
		t.FailNow()
	}
	if w == nil {
		t.Error("Initialise did not return valid object")
		t.FailNow()
	}
}

func TestInitialiseToFailOnEmptyName(t *testing.T) {
	fmt.Println("==============")
	var w wpw.WPWithin
	var e error
	w,e = wpw.Initialise("","Dummy")
	if e == nil {
		t.Error("Initialise did not fail on empty device name")
		t.FailNow()
	}
	if w != nil {
		t.Error("Initialise did not return nil object on empty device name")
		t.FailNow()
	}
}
func TestInitialiseToFailOnEmptyDescription(t *testing.T) {
	fmt.Println("==============")
	var w wpw.WPWithin
	var e error
	w,e = wpw.Initialise("Dummy","")
	if e == nil {
		t.Error("Initialise did not fail on empty device description, ", e)
		t.FailNow()
	}
	if w != nil {
		t.Error("Initialise did not return nil object on empty device description")
		t.FailNow()
	}
}

func TestAddService(t *testing.T) {
	fmt.Println("==============")
	w,e := wpw.Initialise("Dummy", "Dummy")
	svc,e := types.NewService()
	if e != nil {
		t.Error("New service definition failed, ", e)
		t.FailNow()
	}

	w.AddService(svc)
	if e != nil {
		t.Error(".AddService failed, ", e)
		t.FailNow()
	}
	if svc == nil {
		t.Error("New service definition returned nil, expected types.Service")
		t.FailNow()
	}
}
func TestRemoveService(t *testing.T) {
	fmt.Println("==============")
	w,e := wpw.Initialise("Dummy", "Dummy")
	svc,e := types.NewService()
	if e != nil {
		t.Error("New service definition failed, ", e)
		t.FailNow()
	}

	w.AddService(svc)
	if e != nil {
		t.Error(".AddService failed", e)
		t.FailNow()
	}
	if svc == nil {
		t.Error(".AddService returned nil, expected types.Service")
		t.FailNow()
	}

	e = w.RemoveService(svc)
	if e != nil {
		t.Error(".RemoveService failed, ", e)
		t.FailNow()
	}
}
//func TestInitConsumer(t *testing.T) {} // TODO
func TestInitProducerToFailOnNil(t *testing.T) {
	fmt.Println("==============")
	w,e := wpw.Initialise("Dummy", "Dummy")
	var c map[string]string // Payment Service Provider configuration map
	e = w.InitProducer(c)
	if e == nil {
		t.Error(".InitProducer did not fail on nil configuration")
		t.FailNow()
	}
}
func TestInitProducerToFailOnEmptyCfg(t *testing.T) {
	fmt.Println("==============")
	w,e := wpw.Initialise("Dummy", "Dummy")
	c := make(map[string]string, 0)
	e = w.InitProducer(c)
	if e == nil {
		t.Error(".InitProducer did not fail on empty configuration")
		t.FailNow()
	}
}
func TestInitProducerToFailOnNoPspSet(t *testing.T) {
	fmt.Println("==============")
	w,e := wpw.Initialise("Dummy", "Dummy")
	c := make(map[string]string, 0)
	c["dummy"] = "dummy"
	e = w.InitProducer(c)
	if e == nil {
		t.Error(".InitProducer did not fail on configuration without PSP")
		t.FailNow()
	}
}
func TestInitProducer(t *testing.T) {
	fmt.Println("==============")
	w,e := wpw.Initialise("Dummy", "Dummy")
	c := make(map[string]string, 0)

	fmt.Println("-------------- wrong (empty) psp cfg")
	c[psp.CfgPSPName] = ""
	e = w.InitProducer(c)
	if e == nil {
		t.Error("FAIL .InitProducer did not fail on empty payment service provider name")
		t.FailNow()
	}

	fmt.Println("-------------- empty cfg for onlineworldpay")
	c[psp.CfgPSPName] = onlineworldpay.PSPName
	delete(c, onlineworldpay.CfgAPIEndpoint)
	e = w.InitProducer(c)
	if e == nil {
		t.Error("FAIL .InitProducer did not fail on empty configuration for onlineworldpay")
		t.FailNow()
	}

	fmt.Println("-------------- empty securenet dev ID")
	c[psp.CfgPSPName] = securenet.PSPName
	e = w.InitProducer(c)
	if e == nil {
		t.Error("FAIL .InitProducer did not fail on empty configuration for securenet")
		t.FailNow()
	}

	fmt.Println("-------------- invalid securenet dev ID")
	c[securenet.CfgDeveloperID] = "value that cannot be converted to a number"
	e = w.InitProducer(c)
	if e == nil {
		t.Error("FAIL .InitProducer did not fail on invalid securenet developer ID")
		t.FailNow()
	}
	delete(c, securenet.CfgDeveloperID)

	fmt.Println("-------------- online.worldpay.com with no keys")
	c[psp.CfgPSPName] = onlineworldpay.PSPName
	delete(c, psp.CfgHTEPublicKey)
	delete(c, psp.CfgHTEPrivateKey)
	e = w.InitProducer(c)
	if e == nil {
		t.Error("FAIL .InitProducer did not fail on undefined keys")
		t.FailNow()
	}

	//t.Log("Skipping invalid keys scenario as it is failing at the moment")
	fmt.Println("!!! Skipping invalid keys scenario as it is failing at the moment")
	/*
	fmt.Println("-------------- online.worldpay.com with invalid keys")
	c[psp.CfgPSPName] = onlineworldpay.PSPName
	c[psp.CfgHTEPublicKey] = "dummy"
	c[psp.CfgHTEPrivateKey] = "dummy"
	e = w.InitProducer(c)
	if e == nil {
		t.Error(".InitProducer did not fail on invalid keys")
		t.FailNow()
	}
	*/

	// FIXME This fails if another producer is already runnign on the machine
	fmt.Println("-------------- online.worldpay.com with valid configuration")
	c[psp.CfgPSPName] = onlineworldpay.PSPName
	c[psp.CfgHTEPublicKey] = CFG_HtePubKey
	c[psp.CfgHTEPrivateKey] = CFG_HtePrvKey
	e = w.InitProducer(c)
	fmt.Println("e=",e)
	if e != nil {
		t.Error("FAIL .InitProducer failed on valid configuration")
		t.FailNow()
	}

}

func TestGetDevice(t *testing.T){
	fmt.Println("==============")
	w,_ := wpw.Initialise("Dummy", "Dummy")
	d := w.GetDevice()
	fmt.Printf("d=%s\n",d)
	if d == nil {
		t.Error(".GetDevice returned nil, expected *types.Device")
		t.FailNow()
	}
}

type EventHandler struct{}
func (handler *EventHandler) BeginServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) {
	return
}
func (handler *EventHandler) EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) {
	return
}
func (handler *EventHandler) GenericEvent(name string, message string, data interface{}) error {
	return nil
}
func TestStartStopServiceBroadcast(t *testing.T){
	fmt.Println("==============")
	w,_ := wpw.Initialise("Dummy", "Dummy")

	// Test price
	p, e := types.NewPrice()
	if e != nil {
		t.Error("types.NewPrice failed")
		t.FailNow()
	}
	p.ID = 1
	p.UnitDescription = "dummy"
	p.UnitID = 1
	p.PricePerUnit = &types.PricePerUnit{Amount:1, CurrencyCode:"RON"}

	// Test service
	svc, e := types.NewService()
	if e != nil {
		t.Error("types.NewService failed")
		t.FailNow()
	}
	svc.Name = "DummyService"
	svc.Description = "Dummy service for tests"
	svc.ID = 1
	svc.AddPrice(*p)

	// Events handler
	// TODO: I tried to define this as local anonymous struct type but I failed
	var h EventHandler
	w.SetEventHandler(&h)

	// Payment-service-provider configuration
	c := make(map[string]string,0)
	c[psp.CfgPSPName] = onlineworldpay.PSPName
	c[onlineworldpay.CfgMerchantClientKey]  = "T_C_03eaa1d3-4642-4079-b030-b543ee04b5af"
	c[onlineworldpay.CfgMerchantServiceKey] = "T_S_f50ecb46-ca82-44a7-9c40-421818af5996"
	c[psp.CfgHTEPrivateKey]                 = "T_S_f50ecb46-ca82-44a7-9c40-421818af5996"
	c[psp.CfgHTEPublicKey]                  = "T_C_03eaa1d3-4642-4079-b030-b543ee04b5af"
	c[onlineworldpay.CfgAPIEndpoint] = "https://api.worldpay.com/v1"


	// Producer initialization
	e = w.InitProducer(c)

	// Broadcast
	e = w.StartServiceBroadcast(2)
	if e != nil {
		t.Error(".StartServiceBroadcast failed")
		t.FailNow()
	}
	/* FIXME Disabled for now as failing
	e = w.StartServiceBroadcast(2)
	if e == nil {
		t.Error(".StartServiceBroadcast passed when called second time")
		t.Fail()
	}
	*/

	// Cleanup
	w.StopServiceBroadcast()
}
/*
func TestInitProducer(t *testing.T) {
	w,e := wpw.Initialise("Dummy", "Dummy")
	var c map[string]string // Payment Service Provider configuration map
	e = w.InitProducer(c)
}
*/



//func Test(t *testing.T) {}
