package main

import (
	"fmt"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
)

// EventHandlerImpl Handle events from the SDK Core
type EventHandlerImpl struct{}

// BeginServiceDelivery Called when the SDK accepts a call to begin service delivery to a client
func (eh *EventHandlerImpl) BeginServiceDelivery(serviceID int, servicePriceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) {

	fmt.Println("go event from core - begin service delivery")
}

// EndServiceDelivery Called when the SDK accepts a call to end service delivery to a client
func (eh *EventHandlerImpl) EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) {

	fmt.Println("go event from core - end service delivery")
}

// GenericEvent used for handling generic events
func (eh *EventHandlerImpl) GenericEvent(name string, message string, data interface{}) error {

	return nil
}

func (eh *EventHandlerImpl) MakePaymentEvent(totalPrice int, orderCurrency string, clientToken string, orderDescription string, uuid string) {
	return
}

func (eh *EventHandlerImpl) ServiceDiscoveryEvent(remoteAddr string) {
	return
}

func (eh *EventHandlerImpl) ServicePricesEvent(remoteAddr string, serviceId int) {
	return
}

func (eh *EventHandlerImpl) ServiceTotalPriceEvent(remoteAddr string, serviceId int, totalPrice *types.TotalPriceResponse) {
	return
}

func (eh *EventHandlerImpl) ErrorEvent(msg string) {
	return
}
