package types

// Order is placed by consumers. This object holds a record of that order
type Order struct {
	UUID                  string
	Service               Service
	SelectedPriceID       int
	ClientID              string
	SelectedNumberOfUnits int
	PSPReference          string
	DeliveryToken         ServiceDeliveryToken
	ConsumedUnits         int
}
