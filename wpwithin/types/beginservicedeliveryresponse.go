package types

// BeginServiceDeliveryResponse represents a requet to begin service delivery
type BeginServiceDeliveryResponse struct {
	ServerID             string               `json:"serverID"`
	ClientID             string               `json:"clientID"`
	ServiceDeliveryToken ServiceDeliveryToken `json:"serviceDeliveryToken"`
	UnitsToSupply        int                  `json:"unitsToSupply"`
}
