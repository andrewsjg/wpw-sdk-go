package types

// EndServiceDeliveryResponse represents a response to a request to end service delivery
type EndServiceDeliveryResponse struct {
	ServerID             string               `json:"serverID"`
	ClientID             string               `json:"clientID"`
	ServiceDeliveryToken ServiceDeliveryToken `json:"serviceDeliveryToken"`
	UnitsJustSupplied    int                  `json:"unitsJustSupplied"`
	UnitsRemaining       int                  `json:"unitsRemaining"`
}
