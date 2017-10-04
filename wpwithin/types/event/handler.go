package event

import "github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"

// Handler defines events fired by this SDK
type Handler interface {
	BeginServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int)
	EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int)
	GenericEvent(name string, message string, data interface{}) error
	MakePaymentEvent(totalPrice int, orderCurrency string, clientToken string, orderDescription string, uuid string)
	ServiceDiscoveryEvent(remoteAddr string)
	ServicePricesEvent(remoteAddr string, serviceId int)
	ServiceTotalPriceEvent(remoteAddr string, serviceId int, totalPrice *types.TotalPriceResponse)
	ErrorEvent(msg string)
}
