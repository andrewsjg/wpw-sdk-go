package types

// BeginServiceDeliveryRequest represents a request to begin service delivery
type BeginServiceDeliveryRequest struct {
	ClientID             string               `json:"clientID"`
	ServiceDeliveryToken ServiceDeliveryToken `json:"serviceDeliveryToken"`
	UnitsToSupply        int                  `json:"unitsToSupply"`
}
