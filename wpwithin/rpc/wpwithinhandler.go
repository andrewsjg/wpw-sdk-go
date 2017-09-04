package rpc

import (
	log "github.com/sirupsen/logrus"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/rpc/wpthrift/gen-go/wpthrift_types"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types/event"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/utils"
)

// WPWithinHandler handle RPC requests
type WPWithinHandler struct {
	wpwithin wpwithin.WPWithin
	callback event.Handler
}

// NewWPWithinHandler create a new instance of WPWithinHandler
func NewWPWithinHandler(wpWithin wpwithin.WPWithin, callback event.Handler) *WPWithinHandler {

	log.Debug("Begin RPC.WPWithinHandler.NewWPWithinHander()")

	result := &WPWithinHandler{
		wpwithin: wpWithin,
		callback: callback,
	}

	log.Debug("End RPC.WPWithinHandler.NewWPWithinHander()")

	return result
}

// Setup a device with name and description
func (wp *WPWithinHandler) Setup(name, description string) (err error) {

	log.Debug("Begin RPC.WPWithinHandler.Setup()")

	defer log.Debug("End RPC.WPWithinHandler.Setup()")

	wpw, err := wpwithin.Initialise(name, description)

	if err != nil {

		log.Debugf("Error initialising WPWithin. Error = %s", err.Error())

		return err
	}

	wp.wpwithin = wpw

	if wp.callback != nil {

		log.Debug("wp.callback is set, setting handler in WPWithin.")

		wp.wpwithin.SetEventHandler(wp.callback)
	} else {

		log.Debug("wp.callback not set, not setting handler in WPWithin.")
	}

	return nil
}

// AddService Add a new service
func (wp *WPWithinHandler) AddService(svc *wpthrift_types.Service) (err error) {

	log.Debug("Begin RPC.WPWithinHandler.AddService()")

	nSvc := convertThriftServiceToNative(svc)

	log.Debug("End RPC.WPWithinHandler.AddService()")

	return wp.wpwithin.AddService(nSvc)
}

// RemoveService remove a service
func (wp *WPWithinHandler) RemoveService(svc *wpthrift_types.Service) (err error) {

	log.Debug("Begin RPC.WPWithinHandler.RemoveService()")

	nSvc := convertThriftServiceToNative(svc)

	log.Debug("End RPC.WPWithinHandler.RemoveService()")

	return wp.wpwithin.RemoveService(nSvc)
}

// InitConsumer Initialise a consumer to connect to a producer
func (wp *WPWithinHandler) InitConsumer(scheme string, hostname string, port int32, urlPrefix string, serviceID string, hceCard *wpthrift_types.HCECard, pspConfig map[string]string) (err error) {

	log.Debug("RPC.WPWithinHandler.InitConsumer()")

	_hceCard := types.HCECard{
		FirstName:  hceCard.FirstName,
		LastName:   hceCard.LastName,
		ExpMonth:   hceCard.ExpMonth,
		ExpYear:    hceCard.ExpYear,
		CardNumber: hceCard.CardNumber,
		Type:       hceCard.Type,
		Cvc:        hceCard.Cvc,
	}

	return wp.wpwithin.InitConsumer(scheme, hostname, int(port), urlPrefix, serviceID, &_hceCard, pspConfig)
}

// InitProducer initialise a producer
func (wp *WPWithinHandler) InitProducer(pspConfig map[string]string) (err error) {

	log.Debug("RPC.WPWithinHandler.InitProducer()")

	return wp.wpwithin.InitProducer(pspConfig)
}

// GetDevice returns details of the running device
func (wp *WPWithinHandler) GetDevice() (r *wpthrift_types.Device, err error) {

	log.Debug("Begin RPC.WPWithinHandler.GetDevice()")

	device := wp.wpwithin.GetDevice()

	result := &wpthrift_types.Device{

		UID:         device.UID,
		Name:        device.Name,
		Description: device.Description,
		Services:    make(map[int32]*wpthrift_types.Service, 0),
		Ipv4Address: device.IPv4Address,
	}

	log.Debugf("Found %d services for device", len(device.Services))

	if device != nil && len(device.Services) > 0 {

		log.Debug("Begin convert Go Service type to Thrift Service type")

		// Convert the services to Thrift services
		for i, svc := range device.Services {

			// Convert the prices to Thrift prices
			svcPrices := svc.Prices
			thriftPrices := make(map[int32]*wpthrift_types.Price, 0)

			log.Debugf("Found %d prices for service: %s (%d)", len(svcPrices), svc.ID, svc.Name)

			if len(svcPrices) > 0 {

				log.Debug("Begin convert Go price type to Thrift price type")

				for _, svcPrice := range svcPrices {

					thriftPrices[int32(svcPrice.ID)] = &wpthrift_types.Price{

						ID:          int32(svcPrice.ID),
						Description: svcPrice.Description,
						PricePerUnit: &wpthrift_types.PricePerUnit{
							Amount:       int32(svcPrice.PricePerUnit.Amount),
							CurrencyCode: svcPrice.PricePerUnit.CurrencyCode,
						},
						UnitId:          int32(svcPrice.UnitID),
						UnitDescription: svcPrice.UnitDescription,
					}
				}

				log.Debug("End convert Go price type to Thrift price type")
			}

			result.Services[int32(i)] = &wpthrift_types.Service{

				ID:          int32(svc.ID),
				Name:        svc.Name,
				Description: svc.Description,
				Prices:      thriftPrices,
			}
		}

		log.Debug("End convert Go Service type to Thrift Service type")
	}

	log.Debug("End RPC.WPWithinHandler.GetDevice()")

	return result, nil
}

// StartServiceBroadcast starts broadcasting presence of device on network
func (wp *WPWithinHandler) StartServiceBroadcast(timeoutMillis int32) (err error) {

	log.Debug("RPC.WPWithinHandler.StartServiceBroadcast()")

	return wp.wpwithin.StartServiceBroadcast(int(timeoutMillis))
}

// StopServiceBroadcast stops broadcasting presence of device on network
func (wp *WPWithinHandler) StopServiceBroadcast() (err error) {

	log.Debug("Begin RPC.WPWithinHandler.StopServiceBroadcast()")

	wp.wpwithin.StopServiceBroadcast()

	log.Debug("End RPC.WPWithinHandler.StopServiceBroadcast()")

	return nil
}

// DeviceDiscovery initiate a discover process to detect the presence of producer devices on the network
func (wp *WPWithinHandler) DeviceDiscovery(timeoutMillis int32) (r map[*wpthrift_types.ServiceMessage]struct{}, err error) {

	log.Debug("Begin RPC.WPWithinHandler.ServiceDiscovery()")

	gSvcMsgs, err := wp.wpwithin.DeviceDiscovery(int(timeoutMillis))

	if err != nil {

		return nil, err
	}

	result := make(map[*wpthrift_types.ServiceMessage]struct{}, 0)

	for _, gSvcMsg := range gSvcMsgs {

		tmp := &wpthrift_types.ServiceMessage{
			DeviceDescription: gSvcMsg.DeviceDescription,
			Hostname:          gSvcMsg.Hostname,
			PortNumber:        int32(gSvcMsg.PortNumber),
			ServerId:          gSvcMsg.ServerID,
			UrlPrefix:         gSvcMsg.URLPrefix,
			Scheme:            gSvcMsg.Scheme,
			DeviceName:        gSvcMsg.DeviceName,
		}

		result[tmp] = struct{}{}
	}

	log.Debug("End RPC.WPWithinHandler.ServiceDiscovery()")

	return result, nil
}

// RequestServices from a consumers perspective, request the services offered by a consumer
func (wp *WPWithinHandler) RequestServices() (r map[*wpthrift_types.ServiceDetails]struct{}, err error) {

	log.Debug("Begin RPC.WPWithinHandler.RequestServices()")

	gServices, err := wp.wpwithin.RequestServices()

	if err != nil {

		return nil, err
	}

	result := make(map[*wpthrift_types.ServiceDetails]struct{}, 0)

	for _, gService := range gServices {

		tmp := &wpthrift_types.ServiceDetails{
			ServiceId:          int32(gService.ServiceID),
			ServiceDescription: gService.ServiceDescription,
		}

		result[tmp] = struct{}{}
	}

	log.Debug("End RPC.WPWithinHandler.RequestServices()")

	return result, nil
}

// GetServicePrices from a consumers perspective, get the prices for a particular service
func (wp *WPWithinHandler) GetServicePrices(serviceID int32) (r map[*wpthrift_types.Price]struct{}, err error) {

	log.Debug("RPC.WPWithinHandler.GetServicePrices()")

	gSvcPrices, err := wp.wpwithin.GetServicePrices(int(serviceID))

	if err != nil {

		return nil, err
	}

	result := make(map[*wpthrift_types.Price]struct{}, 0)

	for _, gSvcPrice := range gSvcPrices {

		tmp := &wpthrift_types.Price{
			ID:          int32(gSvcPrice.ID),
			Description: gSvcPrice.Description,
			PricePerUnit: &wpthrift_types.PricePerUnit{
				Amount:       int32(gSvcPrice.PricePerUnit.Amount),
				CurrencyCode: gSvcPrice.PricePerUnit.CurrencyCode,
			},
			UnitId:          int32(gSvcPrice.UnitID),
			UnitDescription: gSvcPrice.UnitDescription,
		}

		result[tmp] = struct{}{}
	}

	return result, nil
}

// SelectService from a consumers perspective, select a service to pay for
// Also select the desired price and the number of units
func (wp *WPWithinHandler) SelectService(serviceID int32, numberOfUnits int32, priceID int32) (r *wpthrift_types.TotalPriceResponse, err error) {

	log.Debug("Begin RPC.WPWithinHandler.SelectService()")

	gPriceResponse, err := wp.wpwithin.SelectService(int(serviceID), int(numberOfUnits), int(priceID))

	if err != nil {

		return nil, err
	}

	result := &wpthrift_types.TotalPriceResponse{
		ServerId:           gPriceResponse.ServerID,
		ClientId:           gPriceResponse.ClientID,
		PriceId:            int32(gPriceResponse.PriceID),
		UnitsToSupply:      int32(gPriceResponse.UnitsToSupply),
		TotalPrice:         int32(gPriceResponse.TotalPrice),
		PaymentReferenceId: gPriceResponse.PaymentReferenceID,
		MerchantClientKey:  gPriceResponse.MerchantClientKey,
		CurrencyCode:       gPriceResponse.CurrencyCode,
	}

	log.Debug("End RPC.WPWithinHandler.SelectService()")

	return result, nil
}

// MakePayment a consumer calls this to pay for a selected service
func (wp *WPWithinHandler) MakePayment(request *wpthrift_types.TotalPriceResponse) (r *wpthrift_types.PaymentResponse, err error) {

	log.Debug("Begin RPC.WPWithinHandler.MakePayment()")

	gRequest := types.TotalPriceResponse{
		ServerID:           request.ServerId,
		ClientID:           request.ClientId,
		PriceID:            int(request.PriceId),
		UnitsToSupply:      int(request.UnitsToSupply),
		TotalPrice:         int(request.TotalPrice),
		PaymentReferenceID: request.PaymentReferenceId,
		MerchantClientKey:  request.MerchantClientKey,
		CurrencyCode:       request.CurrencyCode,
	}

	log.Debug("Finised converting thrift.TotalPriceResponse to go.TotalPriceResponse")

	log.Debug("Proceeding to call MakePayment internally using converted request object")

	gPaymentResponse, err := wp.wpwithin.MakePayment(gRequest)

	if err != nil {

		return nil, err
	}

	// TODO create delivery token manually and assign to paymentresponse - need automatpping
	deliveryToken := &wpthrift_types.ServiceDeliveryToken{

		Key:            gPaymentResponse.ServiceDeliveryToken.Key,
		Issued:         utils.TimeFormatISO(gPaymentResponse.ServiceDeliveryToken.Issued),
		Expiry:         utils.TimeFormatISO(gPaymentResponse.ServiceDeliveryToken.Expiry),
		RefundOnExpiry: gPaymentResponse.ServiceDeliveryToken.RefundOnExpiry,
		Signature:      gPaymentResponse.ServiceDeliveryToken.Signature,
	}

	result := &wpthrift_types.PaymentResponse{
		ServerId:             gPaymentResponse.ServerID,
		ClientId:             gPaymentResponse.ClientID,
		TotalPaid:            int32(gPaymentResponse.TotalPaid),
		ServiceDeliveryToken: deliveryToken,
	}

	log.Debug("End RPC.WPWithinHandler.MakePayment()")

	return result, nil
}

// BeginServiceDelivery begin the delivery of a purchased service
func (wp *WPWithinHandler) BeginServiceDelivery(serviceID int32, serviceDeliveryToken *wpthrift_types.ServiceDeliveryToken, unitsToSupply int32) (*wpthrift_types.ServiceDeliveryToken, error) {

	log.WithFields(log.Fields{"serviceID": serviceID, "serviceDeliveryToken": serviceDeliveryToken, "unitsToSupply": unitsToSupply}).Debug("begin rpc.WPWithinHandler.BeginServiceDelivery()")

	defer log.Debug("end rpc.WPWithinHandler.BeginServiceDelivery()")

	issueTime, err := utils.ParseISOTime(serviceDeliveryToken.Issued)

	if err != nil {

		log.Debugf("Error parsing serviceDeliveryToken.Issued time into ISOTime. Error = %s", err.Error())
		return nil, err
	}

	expiryTime, err := utils.ParseISOTime(serviceDeliveryToken.Expiry)

	if err != nil {

		log.Debugf("Error parsing serviceDeliveryToken.Expiry time into ISOTime. Error = %s", err.Error())
		return nil, err
	}

	sdt := types.ServiceDeliveryToken{

		Key:            serviceDeliveryToken.Key,
		Issued:         issueTime,
		Expiry:         expiryTime,
		RefundOnExpiry: serviceDeliveryToken.RefundOnExpiry,
		Signature:      serviceDeliveryToken.Signature,
	}

	_sdt, err := wp.wpwithin.BeginServiceDelivery(int(serviceID), sdt, int(unitsToSupply))

	if err != nil {

		return nil, err
	}

	return convertDeliveryTokenToThrift(_sdt)
}

// EndServiceDelivery end deliery of a purchased service
func (wp *WPWithinHandler) EndServiceDelivery(serviceID int32, serviceDeliveryToken *wpthrift_types.ServiceDeliveryToken, unitsReceived int32) (*wpthrift_types.ServiceDeliveryToken, error) {

	log.WithFields(log.Fields{"serviceID": serviceID, "serviceDeliveryToken": serviceDeliveryToken, "unitsReceived": unitsReceived}).Debug("begin rpc.WPWithinHandler.EndServiceDelivery()")

	defer log.Debug("end rpc.WPWithinHandler.EndServiceDelivery()")

	issueTime, err := utils.ParseISOTime(serviceDeliveryToken.Issued)

	if err != nil {

		return nil, err
	}

	expiryTime, err := utils.ParseISOTime(serviceDeliveryToken.Expiry)

	if err != nil {

		return nil, err
	}

	sdt := types.ServiceDeliveryToken{

		Key:            serviceDeliveryToken.Key,
		Issued:         issueTime,
		Expiry:         expiryTime,
		RefundOnExpiry: serviceDeliveryToken.RefundOnExpiry,
		Signature:      serviceDeliveryToken.Signature,
	}

	_sdt, err := wp.wpwithin.EndServiceDelivery(int(serviceID), sdt, int(unitsReceived))

	if err != nil {

		return nil, err
	}

	return convertDeliveryTokenToThrift(_sdt)
}

// call closeRPCAgent to close RPC
func (wp *WPWithinHandler) CloseRPCAgent() error {

	log.Debug("RPC.WPWithinHandler.closeRPCAgent()")
	wp.CloseRPCAgent()
	return nil
}

func convertThriftServiceToNative(tSvc *wpthrift_types.Service) *types.Service {

	nSvc, _ := types.NewService()

	nSvc.ID = int(tSvc.ID)
	nSvc.Name = tSvc.Name
	nSvc.Description = tSvc.Description

	for _, tPrice := range tSvc.Prices {

		nPrice, _ := types.NewPrice()

		nPrice.Description = tPrice.GetDescription()
		nPrice.ID = int(tPrice.GetID())
		nPrice.UnitDescription = tPrice.GetUnitDescription()
		nPrice.UnitID = int(tPrice.GetUnitId())

		nPpu := &types.PricePerUnit{}

		nPpu.Amount = int(tPrice.GetPricePerUnit().Amount)
		nPpu.CurrencyCode = tPrice.GetPricePerUnit().CurrencyCode

		nPrice.PricePerUnit = nPpu

		nSvc.AddPrice(*nPrice)
	}

	return nSvc
}

func convertDeliveryTokenToThrift(sdt types.ServiceDeliveryToken) (*wpthrift_types.ServiceDeliveryToken, error) {

	tSDT := wpthrift_types.NewServiceDeliveryToken()
	tSDT.Issued = utils.TimeFormatISO(sdt.Issued)
	tSDT.Expiry = utils.TimeFormatISO(sdt.Expiry)
	tSDT.Key = sdt.Key
	tSDT.RefundOnExpiry = sdt.RefundOnExpiry
	tSDT.Signature = sdt.Signature

	return tSDT, nil
}
