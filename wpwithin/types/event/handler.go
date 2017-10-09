package event

import "github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"

// Handler defines events fired by this SDK
type Handler interface {
	BeginServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int)
	EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int)
	GenericEvent(name string, message string, data interface{}) error
}
