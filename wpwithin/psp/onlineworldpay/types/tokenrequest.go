package types

// TokenRequest HTTP Message
type TokenRequest struct {
	Reusable bool `json:"reusable"`

	PaymentMethod TokenRequestPaymentMethod `json:"paymentMethod"`

	ClientKey string `json:"clientKey"`
}
