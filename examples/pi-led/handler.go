package main

import (
	"fmt"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
	"github.com/stianeikeland/go-rpio"
)

// Handler handles the events coming from Worldpay Within
type Handler struct {
	ledGreen  rpio.Pin
	ledRed    rpio.Pin
	ledYellow rpio.Pin
}

func (handler *Handler) setup() error {

	handler.ledGreen = rpio.Pin(2)
	handler.ledRed = rpio.Pin(3)
	handler.ledYellow = rpio.Pin(4)

	if err := rpio.Open(); err != nil {

		return err
	}

	// Cleanup (defer until end)
	// rpio.Close()

	// Ensure pins are in output mode
	handler.ledGreen.Output()
	handler.ledRed.Output()
	handler.ledYellow.Output()

	// Turn of both LEDs, set the pins to low.
	handler.ledGreen.Low()
	handler.ledRed.Low()
	handler.ledYellow.Low()

	return nil
}

// BeginServiceDelivery is called by Worldpay Within when a consumer wish to begin delivery of a service
func (handler *Handler) BeginServiceDelivery(serviceID int, servicePriceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) {

	fmt.Printf("BeginServiceDelivery. ServiceID = %d\n", serviceID)
	fmt.Printf("BeginServiceDelivery. ServicePriceID = %d\n", servicePriceID)
	fmt.Printf("BeginServiceDelivery. UnitsToSupply = %d\n", unitsToSupply)
	fmt.Printf("BeginServiceDelivery. DeliveryToken = %+v\n", serviceDeliveryToken)

	if serviceID == 1 {

		handler.ledGreen.High()

	} else if serviceID == 2 {

		handler.ledRed.High()
	} else if serviceID == 3 {

		handler.ledYellow.High()
	}
}

// EndServiceDelivery is called by Worldpay Within when a consumer wish to end delivery of a service
func (handler *Handler) EndServiceDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) {

	fmt.Printf("EndServiceDelivery. ServiceID = %d\n", serviceID)
	fmt.Printf("EndServiceDelivery. UnitsReceived = %d\n", unitsReceived)
	fmt.Printf("EndServiceDelivery. DeliveryToken = %+v\n", serviceDeliveryToken)

	if serviceID == 1 {

		handler.ledGreen.Low()

	} else if serviceID == 2 {

		handler.ledRed.Low()
	} else if serviceID == 3 {

		handler.ledYellow.Low()
	}
}

// GenericEvent handle a generic event
func (handler *Handler) GenericEvent(name string, message string, data interface{}) error {

	return nil
}

func (handler *Handler) MakePaymentEvent(totalPrice int, orderCurrency string, clientToken string, orderDescription string, uuid string) {
	return
}

func (handler *Handler) ServiceDiscoveryEvent(remoteAddr string) {
	return
}

func (handler *Handler) ServicePricesEvent(remoteAddr string, serviceId int) {
	return
}

func (handler *Handler) ServiceTotalPriceEvent(remoteAddr string, serviceId int, totalPrice *types.TotalPriceResponse) {
	return
}

func (handler *Handler) ErrorEvent(msg string) {
	return
}
