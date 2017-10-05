package rpc

import (
	"crypto/tls"
	"errors"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	log "github.com/sirupsen/logrus"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/rpc/wpthrift/gen-go/wpthrift"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/rpc/wpthrift/gen-go/wpthrift_types"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types/event"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/utils"
)

// EventSenderImpl implementation of event.Handler. Used to send events over Thrift RPC
type EventSenderImpl struct {
	client          wpthrift.WPWithinCallback
	connected       bool
	protocolFactory thrift.TProtocolFactory
	transport       thrift.TTransport
}

// BeginServiceDelivery event
func (cb *EventSenderImpl) BeginServiceDelivery(serviceID int, servicePriceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) {

	log.WithFields(log.Fields{"serviceID": serviceID, "servicePriceID": servicePriceID, "serviceDeliveryToken": fmt.Sprintf("%+v", serviceDeliveryToken), "unitsToSupply": unitsToSupply}).Debug("begin EventSenderImpl.BeginServiceDelivery()")

	defer log.Debug("end EventSenderImpl.BeginServiceDelivery()")

	cb.connectCallbackIfNotConnected()

	sdt := wpthrift_types.ServiceDeliveryToken{

		Key:            serviceDeliveryToken.Key,
		Issued:         utils.TimeFormatISO(serviceDeliveryToken.Issued),
		Expiry:         utils.TimeFormatISO(serviceDeliveryToken.Expiry),
		RefundOnExpiry: serviceDeliveryToken.RefundOnExpiry,
		Signature:      serviceDeliveryToken.Signature,
	}

	err := cb.client.BeginServiceDelivery(int32(serviceID), &sdt, int32(unitsToSupply))

	if err != nil {

		fmt.Println(err.Error())
		log.WithField("error", err.Error()).Error("error calling BeginServiceDelivery using thrift callback client.")
	}
}

// EndServiceDelivery event
func (cb *EventSenderImpl) EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) {

	log.WithFields(log.Fields{"serviceID": serviceID, "serviceDeliveryToken": fmt.Sprintf("%+v", serviceDeliveryToken), "unitsReceived": unitsReceived}).Debug("begin EventSenderImpl.EndServiceDelivery()")

	defer log.Debug("end EventSenderImpl.EndServiceDelivery()")

	cb.connectCallbackIfNotConnected()

	sdt := wpthrift_types.ServiceDeliveryToken{

		Key:            serviceDeliveryToken.Key,
		Issued:         utils.TimeFormatISO(serviceDeliveryToken.Issued),
		Expiry:         utils.TimeFormatISO(serviceDeliveryToken.Expiry),
		RefundOnExpiry: serviceDeliveryToken.RefundOnExpiry,
		Signature:      serviceDeliveryToken.Signature,
	}

	err := cb.client.EndServiceDelivery(int32(serviceID), &sdt, int32(unitsReceived))

	if err != nil {

		fmt.Println(err.Error())
		log.WithField("error", err.Error()).Error("error calling EndServiceDelivery using thrift callback client.")
	}
}

// NewEventSender creates new instance of event sender
func NewEventSender(cfg Configuration) (event.Handler, error) {

	log.WithField("config", fmt.Sprintf("%+v", cfg)).Debug("begin rpc.EventSenderImpl.NewEventSender()")
	defer log.Debug("end rpc.EventSenderImpl.NewEventSender()")

	protocolFactory := thrift.NewTBinaryProtocolFactory(true, true)
	transportFactory := thrift.NewTBufferedTransportFactory(8192)

	var transport thrift.TTransport
	var err error

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.CallbackPort)

	log.Debugf("Will use callback connection string = %s", addr)

	if cfg.Secure {
		tlsCfg := new(tls.Config)
		tlsCfg.InsecureSkipVerify = true
		transport, err = thrift.NewTSSLSocket(addr, tlsCfg)
	} else {
		transport, err = thrift.NewTSocket(addr)
	}
	if err != nil {

		log.Errorf("Error opening socket. Error = %s", err.Error())

		return nil, err
	}
	transport = transportFactory.GetTransport(transport)
	log.Warn("TODO: Transport not going to close..")
	// TODO - How to close this later?
	//	defer transport.Close()
	// if err := transport.Open(); err != nil {
	// 	return nil, err
	// }

	result := &EventSenderImpl{
		connected:       false,
		transport:       transport,
		protocolFactory: protocolFactory,
	}

	return result, nil
}

// connectCallbackIfNotConnected connect the callback client to the server if not already connected
func (cb *EventSenderImpl) connectCallbackIfNotConnected() error {

	log.Debug("begin EventSenderImpl.connectCallbackIfNotConnected()")

	defer log.Debug("end EventSenderImpl.connectCallbackIfNotConnected()")

	log.Debugf("cb.connected: %t", cb.connected)

	if !cb.connected {

		log.Debug("cb is not connected, attempting to connect now.")

		cb.client = wpthrift.NewWPWithinCallbackClientFactory(cb.transport, cb.protocolFactory)

		if err := cb.transport.Open(); err != nil {

			log.Errorf("Cannot connect to callback RPC server.. did you forget to restart this RPC service? Error = %s", err.Error())

			return err
		}

		log.Debug("Did connect to cb.")

		cb.connected = true
	}

	return nil
}

// GenericEvent log a generic event
func (cb *EventSenderImpl) GenericEvent(name string, message string, data interface{}) error {

	return errors.New("eventsender.GenericEvent() is Not implemented")
}

// MakePaymentEvent
func (cb *EventSenderImpl) MakePaymentEvent(totalPrice int, orderCurrency string, clientToken string, orderDescription string, uuid string) {

	log.WithFields(log.Fields{"totalPrice": totalPrice, "orderCurrency": orderCurrency, "clientToken": clientToken, "orderDescription": orderDescription, "uuid": uuid}).Debug("begin EventSenderImpl.MakePaymentEvent()")

	defer log.Debug("end EventSenderImpl.MakePaymentEvent()")

	cb.connectCallbackIfNotConnected()

	err := cb.client.MakePaymentEvent(int32(totalPrice), orderCurrency, clientToken, orderDescription, uuid)

	if err != nil {
		fmt.Println(err.Error())
		log.WithField("error", err.Error()).Error("error calling MakePaymentEvent using thrift callback.")
	}
}

// ServiceDeliveryEvent
func (cb *EventSenderImpl) ServiceDiscoveryEvent(remoteAddr string) {

	log.WithFields(log.Fields{"remoteAddr": remoteAddr}).Debug("begin EventSenderImpl.ServiceDiscoveryEvent()")

	defer log.Debug("end EventSenderImpl.ServiceDiscoveryEvent()")

	cb.connectCallbackIfNotConnected()

	err := cb.client.ServiceDiscoveryEvent(remoteAddr)

	if err != nil {
		fmt.Println(err.Error())
		log.WithField("error", err.Error()).Error("error calling ServiceDiscoveryEvent using thrift callback.")
	}
}

func (cb *EventSenderImpl) ServicePricesEvent(remoteAddr string, serviceId int) {

	log.WithFields(log.Fields{"remoteAddr": remoteAddr}).Debug("begin EventSenderImpl.ServiceDeliveryEvent()")

	defer log.Debug("end EventSenderImpl.ServicePricesEvent()")

	cb.connectCallbackIfNotConnected()

	err := cb.client.ServicePricesEvent(remoteAddr, int32(serviceId))

	if err != nil {
		fmt.Println(err.Error())
		log.WithField("error", err.Error()).Error("error calling ServicePricesEvent using thrift callback.")
	}
}

func (cb *EventSenderImpl) ServiceTotalPriceEvent(remoteAddr string, serviceId int, totalPriceResp *types.TotalPriceResponse) {

	log.WithFields(log.Fields{"remoteAddr": remoteAddr}).Debug("begin EventSenderImpl.ServiceTotalPriceEvent()")

	defer log.Debug("end EventSenderImpl.ServiceTotalPriceEvent()")

	cb.connectCallbackIfNotConnected()

	tpr := wpthrift_types.TotalPriceResponse{
		ServerId:           totalPriceResp.ServerID,
		ClientId:           totalPriceResp.ClientID,
		PriceId:            int32(totalPriceResp.PriceID),
		UnitsToSupply:      int32(totalPriceResp.UnitsToSupply),
		TotalPrice:         int32(totalPriceResp.TotalPrice),
		PaymentReferenceId: totalPriceResp.PaymentReferenceID,
		MerchantClientKey:  totalPriceResp.MerchantClientKey,
		CurrencyCode:       totalPriceResp.CurrencyCode,
	}
	err := cb.client.ServiceTotalPriceEvent(remoteAddr, int32(serviceId), &tpr)

	if err != nil {
		fmt.Println(err.Error())
		log.WithField("error", err.Error()).Error("error calling ServiceTotalPriceEvent using thrift callback.")
	}
}

func (cb *EventSenderImpl) ErrorEvent(msg string) {
	log.WithFields(log.Fields{"msg": msg}).Debug("begin EventSenderImpl.ErrorEvent()")

	defer log.Debug("end EventSenderImpl.ErrorEvent()")

	cb.connectCallbackIfNotConnected()

	err := cb.client.ErrorEvent(msg)

	if err != nil {
		fmt.Println(err.Error())
		log.WithField("error", err.Error()).Error("error calling ErrorEvent using thrift callback.")
	}
}
