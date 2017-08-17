package types

// ServicePriceResponse HTTP Message
type ServicePriceResponse struct {
	ServerID string  `json:"serverID"`
	Prices   []Price `json:"prices"`
}
