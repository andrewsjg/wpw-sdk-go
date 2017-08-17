package types

// PaymentResponse HTTP Message, response to a request to make payment
type PaymentResponse struct {
	ServerID             string                `json:"serverID"`
	ClientID             string                `json:"clientID"`
	TotalPaid            int                   `json:"totalPaid"`
	ServiceDeliveryToken *ServiceDeliveryToken `json:"serviceDeliveryToken"`
}
