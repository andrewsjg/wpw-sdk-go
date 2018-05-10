package rpc_test

import (
	"os"
	"testing"

	//	"github.com/wptechinnovation/wpw-sdk-go/examples/exutils"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/rpc"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/rpc/wpthrift/gen-go/wpthrift_types"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
	//	log "github.com/sirupsen/logrus"
)

var wpw wpwithin.WPWithin
var rpcAgent rpc.WPWithinHandler

//setup

func mySetupFunction() {
	wpw, _ = wpwithin.Initialise("go-go", "gogo-gogo", "")
	rpcAgent = *rpc.NewWPWithinHandler(wpw, nil)
}

func TestBroadcastProcess(t *testing.T) {

	pspConfig := make(map[string]string, 0)
	pspConfig[psp.CfgPSPName] = "worldpayonlinepayments"
	pspConfig[onlineworldpay.CfgMerchantClientKey] = "1234"
	pspConfig[onlineworldpay.CfgMerchantServiceKey] = "5678"
	pspConfig[psp.CfgHTEPrivateKey] = pspConfig[onlineworldpay.CfgMerchantClientKey]
	pspConfig[psp.CfgHTEPublicKey] = pspConfig[onlineworldpay.CfgMerchantServiceKey]
	pspConfig[onlineworldpay.CfgAPIEndpoint] = "/"
	err := rpcAgent.InitProducer(pspConfig)
	if err != nil {
		t.Fail()
		t.Error("Initializing producer has failed with: ", err)
	}

	err = rpcAgent.StartServiceBroadcast(10000)
	if err != nil {
		t.Fail()
		t.Error("Start service broadcast has failed with: ", err)
	}

	_, err = rpcAgent.DeviceDiscovery(5000)
	if err != nil {
		t.Fail()
		t.Error("Device Discovery has failed with: ", err)
	}

	_, err = rpcAgent.SearchForDevice(5000, "go-go")
	if err != nil {
		t.Fail()
		t.Error("Search for device has failed with: ", err)
	}
}

func TestInitConsumer(t *testing.T) {
	hceCard := &wpthrift_types.HCECard{}
	pspConfig := make(map[string]string, 0)
	pspConfig[psp.CfgPSPName] = "worldpayonlinepayments"
	err := rpcAgent.InitConsumer("dummyScheme", "host", 123, "urlPrefix", "serviceID", hceCard, pspConfig)
	if err != nil {
		t.Fail()
		t.Error("Initializing consumer has failed with: ", err)
	}
}

func TestStopServiceBroadcast(t *testing.T) {
	err := rpcAgent.StopServiceBroadcast()
	if err != nil {
		t.Fail()
		t.Error("Stop service broadcast has failed with: ", err)
	}
}

func TestAddService(t *testing.T) {
	srvc := wpthrift_types.NewService()
	err := rpcAgent.AddService(srvc)
	if err != nil {
		t.Fail()
		t.Error("Add service has failed with: ", err)
	}
}

func TestRemoveService(t *testing.T) {
	srvc := wpthrift_types.NewService()
	err := rpcAgent.RemoveService(srvc)
	if err != nil {
		t.Fail()
		t.Error("Remove service has failed with: ", err)
	}
}

func TestSetup(t *testing.T) {
	err := rpcAgent.Setup("dummy", "dummy", "")
	if err != nil {
		t.Fail()
		t.Error("Setup, eventHandler nil failed with: ", err)
	}
	h := &EventHandler{}
	rpcAgent2 := *rpc.NewWPWithinHandler(wpw, h)
	err = rpcAgent2.Setup("dummy", "dummy", "")
	if err != nil {
		t.Fail()
		t.Error("Setup with set eventHandler failed with: ", err)
	}
	err = rpcAgent.Setup("", "", "")
	if err == nil {
		t.Fail()
		t.Error("Setup did not throw error when it should.", err)
	}
}

func TestGetDevice(t *testing.T) {

	service := wpthrift_types.NewService()

	price := wpthrift_types.NewPrice()
	price.PricePerUnit = wpthrift_types.NewPricePerUnit()
	service.Prices = make(map[int32]*wpthrift_types.Price)
	service.Prices[0] = price
	rpcAgent.AddService(service)
	device, err := rpcAgent.GetDevice()
	if device == nil {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}

}

func TestRequestService(t *testing.T) {
	service := wpthrift_types.NewService()
	price := wpthrift_types.NewPrice()
	price.PricePerUnit = wpthrift_types.NewPricePerUnit()
	service.Prices = make(map[int32]*wpthrift_types.Price)
	service.Prices[0] = price
	service.Name = "serviceName"
	service.Description = "serviceDescription"
	service.ID = 1
	service.ServiceType = "serviceType"
	rpcAgent.AddService(service)
	services, err := rpcAgent.RequestServices()
	if services == nil {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}

func TestMain(m *testing.M) {

	mySetupFunction()
	retCode := m.Run()
	os.Exit(retCode)
}

// Dummy handler implementation

type EventHandler struct{}

func (handler *EventHandler) BeginServiceDelivery(serviceID int, servicePriceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) {
	return
}
func (handler *EventHandler) EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) {
	return
}
func (handler *EventHandler) GenericEvent(name string, message string, data interface{}) error {
	return nil
}

func (handler *EventHandler) MakePaymentEvent(totalPrice int, orderCurrency string, clientToken string, orderDescription string, uuid string) {
	return
}
func (handler *EventHandler) ServiceDiscoveryEvent(remoteAddr string) {
	return
}
func (handler *EventHandler) ServicePricesEvent(remoteAddr string, serviceId int) {
	return
}
func (handler *EventHandler) ServiceTotalPriceEvent(remoteAddr string, serviceId int, totalPrice *types.TotalPriceResponse) {
	return
}
func (handler *EventHandler) ErrorEvent(msg string) {
	return
}
