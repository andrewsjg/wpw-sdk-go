package hte

import (
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/types"
)

// PortRangeMin Minimum allowed port
const PortRangeMin = 1

// PortRangeMax Maximum allowed port
const PortRangeMax = 65535

// Client - Allow interaction with a HTE service
type Client interface {

	// Retrieve services available from a HTE service
	DiscoverServices() (types.ServiceListResponse, error)
	// Get the price variants of a particular service
	GetPrices(serviceID int) (types.ServicePriceResponse, error)
	// For a given service, price and quantity, get the price requested by the service
	NegotiatePrice(serviceID, priceID, numberOfUnits int) (types.TotalPriceResponse, error)
	// Given a negotiated price, make a payment for that order
	MakeHtePayment(paymentReferenceID, clientID, clientToken string) (types.PaymentResponse, error)
	// For a purchased service, begin delivery of that product/service
	StartDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) (types.BeginServiceDeliveryResponse, error)
	// For a purchased service, end delivery of that product/service
	EndDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) (types.EndServiceDeliveryResponse, error)
}
