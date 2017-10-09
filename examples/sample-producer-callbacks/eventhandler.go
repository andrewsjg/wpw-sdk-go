package main

import (
	"fmt"

	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
)

// EventHandlerImpl Handle events from the SDK Core
type EventHandlerImpl struct{}

// BeginServiceDelivery Called when the SDK accepts a call to begin service delivery to a client
func (eh *EventHandlerImpl) BeginServiceDelivery(serviceID int, servicePriceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) {

	fmt.Println("go event from core - begin service delivery")
	fmt.Println("selected price id: ", servicePriceID)
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

	fmt.Printf("go event from core - payment: totalPrice=%d, orderCurrency=%s, clientToken=%s, orderDescription=%s, uui=%s\n",
		totalPrice, orderCurrency, clientToken, orderDescription, uuid)
}

func (eh *EventHandlerImpl) ServiceDiscoveryEvent(remoteAddr string) {

	fmt.Printf("go event from core - service dicovery: remoteAddr: %s\n", remoteAddr)
}

func (eh *EventHandlerImpl) ServicePricesEvent(remoteAddr string, serviceId int) {

	fmt.Printf("go event from core - service prices: remoteAddr: %s, serviceId: %d\n", remoteAddr, serviceId)
}

func (eh *EventHandlerImpl) ServiceTotalPriceEvent(remoteAddr string, serviceId int, totalPrice *types.TotalPriceResponse) {

	fmt.Printf("go event from core - service prices: remoteAddr: %s, serviceId: %d\n", remoteAddr, serviceId)
}

func (eh *EventHandlerImpl) ErrorEvent(msg string) {

	fmt.Printf("go event from core - error: %s\n", msg)
}
