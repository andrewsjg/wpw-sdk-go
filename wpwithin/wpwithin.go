package wpwithin

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/configuration"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/core"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/hte"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types/event"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/utils"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/utils/wslog"

	log "github.com/sirupsen/logrus"
)

// Factory to allow easy creation of
var Factory core.SDKFactory

// WPWithin Worldpay Within SDK
type WPWithin interface {
	AddService(service *types.Service) error
	RemoveService(service *types.Service) error
	InitConsumer(scheme, hostname string, portNumber int, urlPrefix, clientID string, hceCard *types.HCECard, pspConfig map[string]string) error
	InitProducer(config map[string]string) error
	GetDevice() *types.Device
	StartServiceBroadcast(timeoutMillis int) error
	StopServiceBroadcast()
	DeviceDiscovery(timeoutMillis int) ([]types.BroadcastMessage, error)
	RequestServices() ([]types.ServiceDetails, error)
	GetServicePrices(serviceID int) ([]types.Price, error)
	SelectService(serviceID, numberOfUnits, priceID int) (types.TotalPriceResponse, error)
	MakePayment(payRequest types.TotalPriceResponse) (types.PaymentResponse, error)
	BeginServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) (types.ServiceDeliveryToken, error)
	EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) (types.ServiceDeliveryToken, error)
	SetEventHandler(handler event.Handler) error
	CloseRPCAgent() error
}

// Initialise Initialise the SDK - Returns an implementation of WPWithin
// Must provide a device name and description
func Initialise(name, description string) (WPWithin, error) {

	log.WithFields(log.Fields{"name": name, "description": description}).Debug("begin wpwithin.Initialise()")

	defer func() {
		if r := recover(); r != nil {

			log.WithFields(log.Fields{"panic_message": r, "name": name, "description": description, "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.Initialise()")
		}
	}()

	// Parameter validation
	log.Debug("Will perform parameter validation")
	if name == "" {

		log.Error("wpwithin.Initialiase() Name parameter is empty")

		return nil, errors.New("name should not be empty")

	} else if description == "" {

		log.Error("wpwithin.Initialiase() Description parameter is empty")

		return nil, errors.New("description should not be empty")
	}
	log.Debug("Parameter validation passed")

	// Start initialisation tasks
	if Factory == nil {

		log.Debug("Will create new core.SDKFactory")

		_Factory, err := core.NewSDKFactory()
		Factory = _Factory

		if err != nil {
			log.WithField("Error", err).Error("Unable to create new core.SDKFactory")
			return nil, fmt.Errorf("Unable to create SDK Factory: %q", err.Error())
		}
		log.Debug("Did create new core.SDKFactory")
	}

	result := &wpWithinImpl{}

	log.Debug("Will create new core.Core")
	core, err := core.NewCore()

	if err != nil {

		log.WithField("Error", err).Error("Error creating new core.Core")

		return result, err
	}
	log.Debug("Did create new core.Core")

	result.core = core

	// Parse configuration
	rawCfg, err1 := configuration.Load("wpwconfig.json")
	wpwConfig := configuration.WPWithin{}
	if err1 != nil {

		log.WithField("Error", err1).Error("Unable to load configuration file wpwconfig.json")
	}

	wpwConfig.ParseConfig(rawCfg)
	core.Configuration = wpwConfig

	log.Debug("Will call doWebSocketLogSetup()")
	doWebSocketLogSetup(core.Configuration)
	log.Debug("After call doWebSocketLogSetup()")

	log.WithFields(log.Fields{"name": name, "description": description}).Debug("Will call Factory.GetDevice()")
	dev, err2 := Factory.GetDevice(name, description)

	if err2 != nil {

		log.WithField("Error", err2).Error("Error calling Factory.GetDevice()")

		return result, err2
	}

	result.core.Device = dev

	log.Debug("Will call Factory.GetOrderManager()")
	om, err := Factory.GetOrderManager()

	if err != nil {

		log.WithField("Error", err).Error("Error calling Factory.GetOrderManager()")

		return result, err
	}
	log.Debug("After call Factory.GetOrderManager()")

	result.core.OrderManager = om

	log.Debug("Will call Factory.GetSvcBroadcaster()")
	bc, err := Factory.GetSvcBroadcaster(result.core.Device.IPv4Address)

	if err != nil {

		log.WithField("Error", err).Error("Error calling Factory.GetSvcBroadcaster()")

		return result, err
	}
	log.Debug("After call Factory.GetSvcBroadcaster()")

	result.core.SvcBroadcaster = bc

	log.Debug("Will call Factory.GetSvcScanner()")
	sc, err := Factory.GetSvcScanner()

	if err != nil {

		log.WithField("Error", err).Error("Error calling Factory.GetSvcScanner()")

		return result, err
	}
	log.Debug("After call Factory.GetSvcScanner()")

	result.core.SvcScanner = sc

	log.Debug("end wpwithin.Initialise()")
	return result, nil
}

type wpWithinImpl struct {
	core *core.Core
}

func (wp *wpWithinImpl) AddService(service *types.Service) error {

	log.Debug("Begin wpwithin.AddService()")

	defer func() {
		if r := recover(); r != nil {

			log.WithFields(log.Fields{"panic_message": r, "service": fmt.Sprintf("%+v", service), "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.AddService()")
		}
	}()

	log.Debugf("Is wp.core.Device.Services map nil? %t", wp.core.Device.Services == nil)
	if wp.core.Device.Services == nil {

		log.Debug("wp.core.Device.Services is nil, will make now: map[int]*types.Service")
		wp.core.Device.Services = make(map[int]*types.Service, 0)
		log.Debug("Did make wp.core.Device.Services: map[int]*types.Service")
	}

	_, exists := wp.core.Device.Services[service.ID]
	log.Debugf("Does wp.core.Device.Services contain service with id %d? %t", service.ID, exists)

	if exists {

		log.Error("Service with id already exists, returning error.")
		return errors.New("Service with that id already exists")
	}

	log.Debug("Adding service to wp.core.Device.Services as it doesn't appear to exist.")
	wp.core.Device.Services[service.ID] = service

	log.Debug("End wpwithin.AddService()")

	return nil
}

func (wp *wpWithinImpl) RemoveService(service *types.Service) error {

	log.Debug("Begin wpwithin.RemoveService()")

	defer func() {
		if r := recover(); r != nil {

			log.WithFields(log.Fields{"panic_message": r, "service": fmt.Sprintf("%+v", service), "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.RemoveService()")
		}
	}()

	if wp.core.Device.Services != nil {

		log.Debugf("Calling delete %d from wp.core.Device.Services", service.ID)
		delete(wp.core.Device.Services, service.ID)
	}

	log.Debug("End wpwithin.RemoveService()")

	return nil
}

func (wp *wpWithinImpl) InitConsumer(scheme, hostname string, portNumber int, urlPrefix, clientID string, hceCard *types.HCECard, pspConfig map[string]string) error {

	mainLogCtx := log.WithFields(log.Fields{"scheme": scheme, "hostname": hostname, "portnumber": portNumber, "urlPrefix": urlPrefix, "clientID": clientID, "hceCard": hceCard.String(), "pspConfig": pspConfig})
	mainLogCtx.Debug("Begin wpwithin.InitConsumer()")

	defer func() {
		if r := recover(); r != nil {

			log.WithFields(log.Fields{"panic_message": r, "scheme": scheme, "hostname": hostname, "port": portNumber,
				"urlPrefix": urlPrefix, "clientID": clientID, "hceCard": fmt.Sprintf("%+v", hceCard), "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.InitConsumer()")
		}
	}()

	log.Debugf("is pspConfig parameter nil? %t", pspConfig == nil)
	if pspConfig == nil {

		log.Error("PSPConfig map is nil, returning.")

		return errors.New("PSPConfig map must be set")
	}

	// Setup PSP as client
	log.Debug("Will call Factory.GetPSPClient()")
	_psp, err := Factory.GetPSPClient(pspConfig)

	if err != nil {

		log.WithField("Error", err).Error("Error calling Factory.GetPSPClient()")
		return err
	}
	log.Debug("After call Factory.GetPSPClient()")

	wp.core.Psp = _psp

	// Set core HCE Card

	wp.core.HCECard = hceCard

	// Setup HTE Client
	log.Debug("Will call Factory.GetHTEClientHTTP()")
	httpHTE, err := Factory.GetHTEClientHTTP()

	if err != nil {

		log.WithField("Error", err).Error("Error calling Factory.GetHTEClientHTTP()")
		return err
	}
	log.Debug("After call to Factory.GetHTEClientHTTP()")

	mainLogCtx.Debug("Will call hte.NewClient")
	client, err := hte.NewClient(scheme, hostname, portNumber, urlPrefix, clientID, httpHTE)

	if err != nil {

		log.WithField("Error", err).Error("Error calling hte.NewClient()")
		return err
	}
	log.Debug("After call hte.NewClient()")

	wp.core.HTEClient = client

	log.Debug("End wpwithin.InitConsumer()")
	return nil
}

func (wp *wpWithinImpl) InitProducer(pspConfig map[string]string) error {

	log.WithField("pspConfig", pspConfig).Debug("Begin wpwithin.InitProducer()")

	defer func() {
		if r := recover(); r != nil {

			log.WithFields(log.Fields{"panic_message": r, "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.InitProducer()")
		}
	}()

	// Parameter validation
	log.Debugf("Is pspConfig parameter nil? %t", pspConfig == nil)
	if pspConfig == nil {

		log.Error("PSPConfig parameter is nil, returning.")
		return errors.New("PSP Config map must me set")
	}

	// Start HTE initialisation tasks
	log.Debug("Will call Factory.GetPSPMerchant()")
	_psp, err := Factory.GetPSPMerchant(pspConfig)

	if err != nil {

		log.WithField("Error", err).Error("Error calling Factory.GetPSPMerchant()")
		return fmt.Errorf("Unable to create psp: %q", err.Error())
	}
	log.Debug("Did call Factory.GetPSPMerchant()")

	wp.core.Psp = _psp

	log.Debug("Will call hte.NewHTECredential()")
	hteCredential, err := hte.NewHTECredential(pspConfig[psp.CfgHTEPublicKey], pspConfig[psp.CfgHTEPrivateKey])

	if err != nil {

		log.WithField("Error", err).Debug("Error calling hte.NewHTECredential()")
		return err
	}
	log.Debug("Did call hte.NewHTECredential()")

	log.Debug("Will call Factory.GetHTEServiceHandler()")
	hteSvcHandler := Factory.GetHTEServiceHandler(wp.core.Device, wp.core.Psp, hteCredential, wp.core.OrderManager, wp.core.EventHandler)
	log.Debug("Did call Factory.GetHTEServiceHandler()")

	log.Debug("Will call Factory.GetHTE()")
	svc, err := Factory.GetHTE(wp.core.Device, wp.core.Psp, wp.core.Device.IPv4Address, "http://", hteCredential, wp.core.OrderManager, hteSvcHandler)

	if err != nil {

		log.WithField("Error", err).Error("Error calling Factory.GetHTE()")
		return err
	}
	log.Debug("Did call Factory.GetHTE()")

	wp.core.HTE = svc

	log.Debug("Will setup channel of type error make(chan error)")
	// Error channel allows us to get the error out of the go routine
	chStartResult := make(chan error)
	var startErr error
	log.Debug("Did call make(chan error)")
	log.Debug("Will spin up goroutine to start HTE service")
	go func() {

		log.Debug("Inside go routine, will call wp.core.HTE.Start()")
		chStartResult <- wp.core.HTE.Start()
		log.Debug("Inside go routine, did call wp.core.HTE.Start()")
	}()
	log.Debug("Did spin up go routine.. this is main thread.")

	// Receive the error from the channel or wait a predefined amount of time
	// TODO CH : Fix this race condition - Matthew B has a solution, find and implement.
	select {

	case res := <-chStartResult:
		log.Debug("Error channel from go routine contains an error")
		startErr = res

	case <-time.After(time.Millisecond * 750):
		log.Debug("Did wait 750ms, no error")
	}

	log.WithField("Error result", startErr).Debug("End wpwithin.InitProducer()")

	return startErr
}

func (wp *wpWithinImpl) GetDevice() *types.Device {

	log.Debug("Begin wpwithin.GetDevice()")

	defer func() {
		if r := recover(); r != nil {

			log.WithField("Stack", string(debug.Stack())).Errorf("Recover: WPWithin.GetDevice()")
		}
	}()

	result := wp.core.Device
	log.WithField("Device", result).Debug("End wpwithin.GetDevice()")

	return result
}

func (wp *wpWithinImpl) StartServiceBroadcast(timeoutMillis int) error {

	log.Debug("Begin wpwithin.startServiceBroadcast()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Print(string(debug.Stack()))

			log.WithFields(log.Fields{"panic_message": r, "timeoutMillis": timeoutMillis, "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.StartServiceBroadcast()")
		}
	}()

	log.Debug("Will construct types.BroadcastMessage now")
	// Setup message that is broadcast over network
	msg := types.BroadcastMessage{

		DeviceDescription: wp.core.Device.Description,
		Hostname:          wp.core.HTE.IPAddr(),
		ServerID:          wp.core.Device.UID,
		URLPrefix:         wp.core.HTE.URLPrefix(),
		PortNumber:        wp.core.HTE.Port(),
		Scheme:            wp.core.HTE.Scheme(),
		DeviceName:        wp.core.Device.Name,
	}
	log.Debug("Did construct types.BroadcastMessage")

	log.Debug("Will make channel of type error for catching broadcast errors from go routine.")
	// Set up a channel to get the error out of the go routine
	chBroadcastErr := make(chan error)
	var errBroadcast error
	log.Debug("Did make channel of type error")

	log.Debug("Will spin up go routine to perform broadcasting")
	go func() {
		log.WithFields(log.Fields{"message": msg, "timeoutmillis": timeoutMillis}).Debug("Inside go routine, will call wp.core.SvcBroadcaster.StartBroadcast(")
		chBroadcastErr <- wp.core.SvcBroadcaster.StartBroadcast(msg, timeoutMillis)
		log.Debug("Did call wp.core.SvcBroadcaster.StartBroadcast()")
	}()
	log.Debug("Did spin up go routine to perform broadcasting")

	// Either get the error or wait a small amount of time to give the all clear.
	// This is a race condition - ahhhh! TODO CH : Fix this
	select {

	case res := <-chBroadcastErr:

		errBroadcast = res

	case <-time.After(time.Millisecond * 750):

	}

	log.WithField("Result err", errBroadcast).Debug("End wpwithin.startServiceBroadcast()")
	return errBroadcast
}

func (wp *wpWithinImpl) StopServiceBroadcast() {

	log.Debug("Begin wpwithin.StopServiceBroadcast()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithField("Stack", string(debug.Stack())).Errorf("Recover: WPWithin.StopServiceBroadcast()")
		}
	}()

	log.Debug("Will call wp.core.SvcBroadcaster.StopBroadcast()")
	err := wp.core.SvcBroadcaster.StopBroadcast()

	if err != nil {

		log.WithField("Error", err).Error("Error calling wp.core.SvcBroadcaster.StopBroadcast()")
	}

	log.Debug("Did call wp.core.SvcBroadcaster.StopBroadcast()")

	log.Debug("End wpwithin.StopServiceBroadcast()")
}

func (wp *wpWithinImpl) DeviceDiscovery(timeoutMillis int) ([]types.BroadcastMessage, error) {

	log.WithField("timeoutMillis", timeoutMillis).Debug("Begin wpwithin.DeviceDiscovery()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithFields(log.Fields{"panic_message": r, "timeoutMillis": timeoutMillis, "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.DeviceDiscovery()")
		}
	}()

	var svcResults []types.BroadcastMessage

	log.Debug("Will call wp.core.SvcScanner.ScanForServices()")
	if scanResult, err := wp.core.SvcScanner.ScanForServices(timeoutMillis); err != nil {

		log.WithField("Error", err).Error("Error calling wp.core.SvcScanner.ScanForServices()")

		return nil, err

	} else if len(scanResult) > 0 {

		log.WithField("Count devices", len(scanResult)).Debug("Did call wp.core.SvcScanner.ScanForServices()")
		log.Info("Did find %d devices", len(scanResult))

		// Convert map of services to array
		for _, svc := range scanResult {

			svcResults = append(svcResults, svc)
		}
	}
	log.Debug("After call wp.core.SvcScanner.ScanForService()")

	log.WithField("Search result", svcResults).Debug("End wpwithin.DeviceDiscovery()")

	return svcResults, nil
}

func (wp *wpWithinImpl) GetServicePrices(serviceID int) ([]types.Price, error) {

	log.WithField("ServiceID", serviceID).Debug("Begin wpwithin.GetServicePrices()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithFields(log.Fields{"panic_message": r, "serviceID": serviceID, "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.GetServicePrices()")
		}
	}()

	var result []types.Price

	log.WithField("ServiceID", serviceID).Debug("Will call wp.core.HTEClient.GetPrices()")
	priceResponse, err := wp.core.HTEClient.GetPrices(serviceID)

	if err != nil {

		log.WithField("Error", err).Error("Error calling wp.core.HTEClient.GetPrices()")
		return nil, err
	}
	log.Debug("Did call wp.core.HTEClient.GetPrices()")

	if priceResponse.Prices == nil {

		log.Debug("Seemed to be successful call to wp.core.HTEClient.GetPrices() but response.Prices is nil")

		return nil, errors.New("priceResponse.Prices is nil")
	}

	log.Infof("Did find %d prices for service id %d", len(priceResponse.Prices), serviceID)

	for _, price := range priceResponse.Prices {

		result = append(result, price)
	}

	log.WithField("result", result).Debug("End wpwithin.GetServicePrices()")

	return result, nil
}

func (wp *wpWithinImpl) SelectService(serviceID, numberOfUnits, priceID int) (types.TotalPriceResponse, error) {

	log.WithFields(log.Fields{"serviceID": serviceID, "numberOfUnits": numberOfUnits, "priceID": priceID}).Debug("Begin wpwithin.SelectService()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithFields(log.Fields{"panic_message": r, "serviceID": serviceID, "numberOfUnits": numberOfUnits, "priceID": priceID, "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.SelectService()")
		}
	}()

	log.Debug("Will call wp.core.HTEClient.NegotiatePrice()")
	tpr, err := wp.core.HTEClient.NegotiatePrice(serviceID, priceID, numberOfUnits)
	log.Debug("Did call wp.core.HTEClient.NegotiatePrice()")
	if err != nil {

		log.WithField("Error", err).Error("Error calling wp.core.HTEClient.NegotiatePrice()")
	}

	log.WithFields(log.Fields{"Error": err, "totalPriceResponse": tpr}).Debug("End wpwithin.SelectService()")

	return tpr, err
}

func (wp *wpWithinImpl) MakePayment(request types.TotalPriceResponse) (types.PaymentResponse, error) {

	log.WithField("TotalPriceResponse parameter", request).Debug("Begin wpwithin.MakePayment()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithFields(log.Fields{"panic_message": r, "price request": fmt.Sprintf("%+v", request), "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.MakePayment()")
		}
	}()

	log.WithFields(log.Fields{"HCE Card": wp.core.HCECard.String(), "Merchant client key": request.MerchantClientKey}).Debug("Will call wp.core.Psp.GetToken()")
	token, err := wp.core.Psp.GetToken(wp.core.HCECard, request.MerchantClientKey, false)

	if err != nil {

		log.WithField("Error", err).Error("Error calling wp.core.Psp.GetToken()")

		return types.PaymentResponse{}, err
	}
	log.WithField("Token", token).Debug("Did call wp.core.Psp.GetToken() without error.")

	log.WithFields(log.Fields{"Payment reference ID": request.PaymentReferenceID, "Client ID": request.ClientID, "token": token}).Debug("will call wp.core.HTEClient.MakeHtePayment()")
	paymentResponse, err := wp.core.HTEClient.MakeHtePayment(request.PaymentReferenceID, request.ClientID, token)

	if err != nil {

		log.WithField("Error", err).Error("Error calling wp.core.HTEClient.MakeHtePayment()")

	} else {

		log.Debug("Did call wp.core.HTEClient.MakeHtePayment() without error")
	}

	log.WithField("PaymentResponse", paymentResponse).Debug("End wpwithin.MakePayment()")

	return paymentResponse, err
}

func (wp *wpWithinImpl) RequestServices() ([]types.ServiceDetails, error) {

	log.Debug("begin wpwithin.RequestServices()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithField("Stack", string(debug.Stack())).Errorf("Recover: WPWithin.RequestServices()")
		}
	}()

	var result []types.ServiceDetails

	log.Debug("Will call wp.core.HTEClient.DiscoverService()")
	serviceResponse, err := wp.core.HTEClient.DiscoverServices()

	if err != nil {

		log.WithField("Error", err).Error("Error calling wp.core.HTEClient.DiscoverService()")
		return nil, err
	}
	log.Debug("Did call wp.core.HTEClient.DiscoverService()")

	if &serviceResponse == nil {

		log.Error("Call to wp.core.HTEClient.DiscoverServices() seems successful but serviceResponse is nil")
		return nil, errors.New("Call to wp.core.HTEClient.DiscoverServices() seems successful but serviceResponse is nil")
	}

	log.Infof("Did find %d services.", len(serviceResponse.Services))

	for _, svc := range serviceResponse.Services {

		result = append(result, svc)
	}

	log.WithField("Result", result).Debug("end wpwithin.RequestServices()")

	return result, nil
}

func (wp *wpWithinImpl) Core() (*core.Core, error) {

	log.Debug("Begin call to wpwithin.Core() -- just a getter.")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithField("stack", string(debug.Stack())).Errorf("Recover: WPWithin.Core()")
		}
	}()

	log.Debug("End call to wpwithin.Core()")

	return wp.core, nil
}

func (wp *wpWithinImpl) BeginServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) (types.ServiceDeliveryToken, error) {

	log.WithFields(log.Fields{"serviceID": serviceID, "serviceDeliveryToken": fmt.Sprintf("%+v", serviceDeliveryToken), "unitsToSupply": unitsToSupply}).Debug("begin wpwithin.wpwithinimpl.BeginServiceDelivery()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithFields(log.Fields{"panic_message": r, "serviceID": serviceID, "unitsToSupply": unitsToSupply,
				"serviceDeliveryToken": fmt.Sprintf("%+v", serviceDeliveryToken), "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.BeginServiceDelivery()")
		}
	}()

	log.Debug("Will call wp.core.HTEClient.StartDelivery()")
	deliveryResponse, err := wp.core.HTEClient.StartDelivery(serviceID, serviceDeliveryToken, unitsToSupply)

	if err != nil {

		log.Errorf("Error calling beginServiceDelivery. Error: %s", err.Error())
		return types.ServiceDeliveryToken{}, err
	}
	log.WithField("delivery response", deliveryResponse).Debug("Did call wp.core.HTEClient.StartDelivery()")

	log.Debug("end wpwithin.wpwithinimpl.BeginServiceDelivery()")

	return deliveryResponse.ServiceDeliveryToken, nil
}

func (wp *wpWithinImpl) EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) (types.ServiceDeliveryToken, error) {

	log.WithFields(log.Fields{"serviceID": serviceID, "serviceDeliveryToken": fmt.Sprintf("%+v", serviceDeliveryToken), "unitsReceived": unitsReceived}).Debug("begin wpwithin.wpwithinimpl.EndServiceDelivery()")

	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s", debug.Stack())

			log.WithFields(log.Fields{"panic_message": r, "serviceID": serviceID, "unitsReceived": unitsReceived,
				"serviceDeliveryToken": fmt.Sprintf("%+v", serviceDeliveryToken), "stack": fmt.Sprintf("%s", debug.Stack())}).
				Errorf("Recover: WPWithin.EndServiceDelivery()")
		}
	}()

	log.Debug("Will call wp.core.HTEClient.EndDelivery()")
	deliveryResponse, err := wp.core.HTEClient.EndDelivery(serviceID, serviceDeliveryToken, unitsReceived)

	if err != nil {

		log.WithField("Error", err).Error("Error calling wp.core.HTEClient.EndServiceDelivery()")

		return types.ServiceDeliveryToken{}, err
	}
	log.WithField("DeliveryResponse", deliveryResponse).Debug("Did call wp.core.HTEClient.EndDelivery()")

	log.Debug("end wpwithin.wpwithinimpl.EndServiceDelivery()")

	return deliveryResponse.ServiceDeliveryToken, nil
}

func (wp *wpWithinImpl) SetEventHandler(handler event.Handler) error {

	log.Debug("wpwithin.wpwithinimpl setting core event handler")

	wp.core.EventHandler = handler

	return nil
}

func (wp *wpWithinImpl) CloseRPCAgent() error {

	log.Debug("wpwithin.wpwithinimpl calling CloseRPCAgent()")

	os.Exit(0)
	return nil
}

func doWebSocketLogSetup(cfg configuration.WPWithin) {

	log.WithField("Configuration", cfg).Debug("Begin wpwithin.doWebSocketLogSetup()")

	if cfg.WSLogEnable {

		log.Debug("Config.WSLogEnable is true -- proceed to setup WS Logger")

		// Clean up the levels config value - just in case.
		log.Debug("Clean up Config.WSLogLevel -- remove spaces")
		strLevels := strings.Replace(cfg.WSLogLevel, " ", "", -1)
		log.Debug("Split up Config.WSLogLevel to array -- comma delimiter")
		logLevels := strings.Split(strLevels, ",")

		// Support all levels
		var levels []log.Level
		log.WithField("log levels", logLevels).Debug("Setting up logger to support the specified log levels.")
		for _, level := range logLevels {

			switch strings.ToLower(level) {

			case "panic":
				levels = append(levels, log.PanicLevel)
			case "fatal":
				levels = append(levels, log.FatalLevel)
			case "error":
				levels = append(levels, log.ErrorLevel)
			case "warn":
				levels = append(levels, log.WarnLevel)
			case "info":
				levels = append(levels, log.InfoLevel)
			case "debug":
				levels = append(levels, log.DebugLevel)
			}
		}

		log.Debug("Attempt to get external IPv4 address")
		ip, err := utils.ExternalIPv4()
		strIP := ""

		if err == nil {

			log.WithField("IPv4", ip.String()).Debug("Did get external IPv4")
			strIP = ip.String()
		} else {

			log.WithField("Error", err).Error("Failed to get external IPv4 address")
			fmt.Printf("Error getting ExternalIPv4: %s\n", err.Error())
		}

		log.WithFields(log.Fields{"ip_add": strIP, "ws_port": cfg.WSLogPort, "log levels": levels}).Debug("will call wslog.Initialise()")
		err = wslog.Initialise(strIP, cfg.WSLogPort, levels)

		if err != nil {

			log.WithField("Error", err).Error("Error calling wslog.Initialise()")
			fmt.Printf("Error initialising WebSocket logger: %s\n", err.Error())
		} else {

			log.Debug("Did call wslog.Initialise()")
		}
	} else {

		log.Debug("Config.WSLogEnable is false -- didn't setup WS Logger")
	}

	log.Debug("End doWebSocketLogSetup()")
}
