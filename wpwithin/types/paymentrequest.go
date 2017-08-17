package types

// PaymentRequest request to make payment
type PaymentRequest struct {
	ClientID           string `json:"clientID"`
	ClientToken        string `json:"clientToken"`
	PaymentReferenceID string `json:"paymentReferenceID"`
}
