package types

// EndServiceDeliveryRequest represents a request to end service delivery
type EndServiceDeliveryRequest struct {
	ClientID             string               `json:"clientID"`
	ServiceDeliveryToken ServiceDeliveryToken `json:"serviceDeliveryToken"`
	UnitsReceived        int                  `json:"unitsReceived"`
}
