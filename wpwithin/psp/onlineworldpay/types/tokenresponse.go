package types

// TokenResponse HTTP Response Message
type TokenResponse struct {
	Reusable bool `json:"reusable"`

	Token string `json:"token"`

	TokenResponsePaymentMethod TokenResponsePaymentMethod `json:"paymentMethod"`
}
