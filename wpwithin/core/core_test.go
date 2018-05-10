package core

import (
	"testing"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/hte"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/servicediscovery"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
)

var device types.Device

func initializeDevice() {

	service1 := &types.Service{
		ID:          1,
		Name:        "Car charger 1",
		Description: "Can charge your hybrid car",
		Prices: map[int]types.Price{
			1: {
				ID:              1,
				Description:     "low price",
				PricePerUnit:    &types.PricePerUnit{Amount: 25, CurrencyCode: "GBP"},
				UnitID:          1,
				UnitDescription: "One kilowatt-hour",
			},
		},
	}

	service2 := &types.Service{
		ID:          2,
		Name:        "Car charger",
		Description: "Can charge your electric car",
		Prices: map[int]types.Price{
			1: {
				ID:              2,
				Description:     "medium price",
				PricePerUnit:    &types.PricePerUnit{Amount: 50, CurrencyCode: "GBP"},
				UnitID:          2,
				UnitDescription: "One kilowatt per hour",
			},
		},
	}

	var services map[int]*types.Service
	services = map[int]*types.Service{
		1: service1,
		2: service2,
	}

	device = types.Device{
		UID:         "405aae3a-9ed1-11e7-abc4-cec278b6b50a",
		Name:        "John Doe",
		Description: "some description",
		Services:    services,
		IPv4Address: "192.168.35.7",
	}
}

func TestNewCore(t *testing.T) {

	core, err := NewCore()
	if err != nil {
		t.Error("Error should be nil.")
		t.FailNow()
	}
	if core == nil {
		t.Error("Core shouldn't be nil.")
		t.FailNow()
	}
}

func TestSetDevice(t *testing.T) {

	initializeDevice()
	core, _ := NewCore()
	core.SetDevice(&device)

	if core.Device.UID != "405aae3a-9ed1-11e7-abc4-cec278b6b50a" {
		t.Error("UIDs not match.")
		t.FailNow()
	}
}

// dummy implemantation of PSP required for next test
type dummyPSP struct {
	dummy string
}

func (d *dummyPSP) GetToken(hceCredentials *types.HCECard, clientKey string, reusableToken bool) (string, error) {
	return "tokenDummyString", nil
}
func (d *dummyPSP) MakePayment(amount int, currencyCode, clientToken, orderDescription, customerOrderCode string) (string, error) {
	return "paymantDummyString", nil
}
func newDummyImplOfPSP() psp.PSP {
	result := &dummyPSP{
		dummy: "abc123",
	}
	return result
}

func TestSetPsp(t *testing.T) {

	var psp0, psp1 psp.PSP
	psp0 = newDummyImplOfPSP()
	psp1 = newDummyImplOfPSP()
	core, _ := NewCore()
	core.SetPsp(psp0)

	if core.Psp != psp0 {
		t.Error("psp0 should match core.Psp.")
		t.FailNow()
	}
	if core.Psp == psp1 {
		t.Error("psp1 shouldn't match core.Psp.")
		t.FailNow()
	}
}

// dummy implemantation of Broadcaster required for next test
type dummyBroadcaster struct {
	dummyBroadcastRunning bool
}

func (d *dummyBroadcaster) StartBroadcast(msg types.BroadcastMessage, timeoutMillis int) error {
	return nil
}
func (d *dummyBroadcaster) StopBroadcast() error {
	return nil
}

func newDummyImplOfBroadcaster() servicediscovery.Broadcaster {
	result := &dummyBroadcaster{
		dummyBroadcastRunning: false,
	}
	return result
}

func TestSetSvcBroadcaster(t *testing.T) {
	var broadcaster0, broadcaster1 servicediscovery.Broadcaster
	broadcaster0 = newDummyImplOfBroadcaster()
	broadcaster1 = newDummyImplOfBroadcaster()
	core, _ := NewCore()
	core.SetSvcBroadcaster(broadcaster0)
	if core.SvcBroadcaster != broadcaster0 {
		t.Error("broadcaster0 should match core.SvcBroadcaster.")
		t.FailNow()
	}
	if core.SvcBroadcaster == broadcaster1 {
		t.Error("broadcaster1 shouldn't match core.SvcBroadcaster.")
		t.FailNow()
	}
}

// dummy implemantation of Scanner required for next test
type dummyScanner struct {
	dummyScannerRunning bool
}

func (d *dummyScanner) ScanForServices(timeout int) (map[string]types.BroadcastMessage, error) {
	d.dummyScannerRunning = true
	return make(map[string]types.BroadcastMessage), nil
}

func (d *dummyScanner) ScanForService(timeout int, deviceName string) (*types.BroadcastMessage, error) {
	dummyBrMsg := types.BroadcastMessage{}
	return &dummyBrMsg, nil
}

func (d *dummyScanner) StopScanner() {
	d.dummyScannerRunning = false
	return
}

func newDummyImplOfScanner() servicediscovery.Scanner {
	result := &dummyScanner{
		dummyScannerRunning: false,
	}
	return result
}

func TestSetSvcScanner(t *testing.T) {
	var scanner0, scanner1 servicediscovery.Scanner
	scanner0 = newDummyImplOfScanner()
	scanner1 = newDummyImplOfScanner()

	core, _ := NewCore()
	core.SetSvcScanner(scanner0)

	if core.SvcScanner != scanner0 {
		t.Error("scanner0 should match core.SvcScanner.")
		t.FailNow()
	}
	if core.SvcScanner == scanner1 {
		t.Error("scanner1 shouldn't match core.SvcScanner.")
		t.FailNow()
	}
}

// dummy implemantation of Service required for next test
type dummyServiceImpl struct {
	dummyService bool
}

func (d *dummyServiceImpl) Start() error {
	return nil
}
func (d *dummyServiceImpl) IPAddr() string {
	return "192.168.1.2"
}
func (d *dummyServiceImpl) Port() int {
	return 2345
}
func (d *dummyServiceImpl) URLPrefix() string {
	return "http://"
}
func (d *dummyServiceImpl) Scheme() string {
	return "some scheme"
}

func newDummyImplOfService() hte.Service {
	result := &dummyServiceImpl{
		dummyService: false,
	}
	return result
}

func TestSetHTE(t *testing.T) {
	var service0, service1 hte.Service
	service0 = newDummyImplOfService()
	service1 = newDummyImplOfService()

	core, _ := NewCore()
	core.SetHTE(service0)

	if core.HTE != service0 {
		t.Error("service0 should match core.HTE.")
		t.FailNow()
	}
	if core.HTE == service1 {
		t.Error("service1 shouldn't match core.HTE.")
		t.FailNow()
	}
}

func TestSetHCECard(t *testing.T) {
	var hceCard0, hceCard1 types.HCECard
	hceCard0.FirstName = "Joe"
	hceCard1.FirstName = "John"
	core, _ := NewCore()
	core.SetHCECard(&hceCard0)

	if core.HCECard.FirstName != hceCard0.FirstName {
		t.Error("hceCard0.FirstName should match core.HCECard.FirstName.")
		t.FailNow()
	}
	if core.HCECard.FirstName == hceCard1.FirstName {
		t.Error("hceCard1.FirstName shouldn't match core.HCECard.FirstName.")
		t.FailNow()
	}
}

// dummy implemantation of OrderManager required for next test
type dummyOrderManagerImpl struct {
	dummyOrderManagerAttr bool
}

func (d *dummyOrderManagerImpl) AddOrder(order types.Order) error {
	d.dummyOrderManagerAttr = true
	return nil
}
func (d *dummyOrderManagerImpl) GetOrder(orderUUID string) (*types.Order, error) {
	return nil, nil
}
func (d *dummyOrderManagerImpl) OrderExists(orderUUID string) bool {
	return d.dummyOrderManagerAttr
}
func (d *dummyOrderManagerImpl) UpdateOrder(order types.Order) error {
	return nil
}
func newDummyImplOfOrderManager() hte.OrderManager {
	result := &dummyOrderManagerImpl{
		dummyOrderManagerAttr: false,
	}
	return result
}

func TestSetOrderManager(t *testing.T) {
	var orderManager0, orderManager1 hte.OrderManager
	var dummyOrder types.Order

	orderManager0 = newDummyImplOfOrderManager()
	orderManager1 = newDummyImplOfOrderManager()
	orderManager0.AddOrder(dummyOrder) // will change the dummyOrderManagerAttr

	core, _ := NewCore()
	core.SetOrderManager(orderManager0)

	if core.OrderManager.OrderExists("") != orderManager0.OrderExists("") {
		t.Error("core.OrderManager.OrderExists() should equal to orderManager0.OrderExists().")
		t.FailNow()
	}
	if core.OrderManager.OrderExists("") == orderManager1.OrderExists("") {
		t.Error("core.OrderManager.OrderExists() should not equal to orderManager1.OrderExists().")
		t.FailNow()
	}
}

// dummy implemantation of Client required for next test
type dummyClientImpl struct {
	dummyClientAttr bool
}

func (d *dummyClientImpl) DiscoverServices() (types.ServiceListResponse, error) {
	var slr types.ServiceListResponse
	return slr, nil
}
func (d *dummyClientImpl) GetPrices(serviceID int) (types.ServicePriceResponse, error) {
	var spr types.ServicePriceResponse
	return spr, nil
}
func (d *dummyClientImpl) NegotiatePrice(serviceID, priceID, numberOfUnits int) (types.TotalPriceResponse, error) {
	var tpr types.TotalPriceResponse
	return tpr, nil
}
func (d *dummyClientImpl) MakeHtePayment(paymentReferenceID, clientID, clientToken string) (types.PaymentResponse, error) {
	var pr types.PaymentResponse
	return pr, nil
}
func (d *dummyClientImpl) StartDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) (types.BeginServiceDeliveryResponse, error) {
	var bdr types.BeginServiceDeliveryResponse
	return bdr, nil
}
func (d *dummyClientImpl) EndDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) (types.EndServiceDeliveryResponse, error) {
	var esdr types.EndServiceDeliveryResponse
	return esdr, nil
}

func newDummyImplOfClient() hte.Client {
	result := &dummyClientImpl{
		dummyClientAttr: false,
	}
	return result
}

func TestSetHTEClient(t *testing.T) {

	var client0, client1 hte.Client

	client0 = newDummyImplOfClient()
	client1 = newDummyImplOfClient()

	core, _ := NewCore()
	core.SetHTEClient(client0)

	if core.HTEClient != client0 {
		t.Error("core.HTEClient should be equal to client0.")
		t.FailNow()
	}
	if core.HTEClient == client1 {
		t.Error("core.HTEClient should not be equal to client0.")
		t.FailNow()
	}

}
