package main

import (
	"fmt"

	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
)

// EventHandlerImpl Handle events from the SDK Core
type EventHandlerImpl struct{}

// BeginServiceDelivery Called when the SDK accepts a call to begin service delivery to a client
func (eh *EventHandlerImpl) BeginServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) {

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
