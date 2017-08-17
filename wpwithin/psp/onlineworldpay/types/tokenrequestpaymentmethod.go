package types

// TokenRequestPaymentMethod HTTP Message
type TokenRequestPaymentMethod struct {
	Name string `json:"name"`

	ExpiryMonth int32 `json:"expiryMonth"`

	ExpiryYear int32 `json:"expiryYear"`

	IssueNumber *int32 `json:"issueNumber"`

	StartMonth *int32 `json:"startMonth"`

	StartYear *int32 `json:"startYear"`

	CardNumber string `json:"cardNumber"`

	Type string `json:"type"`

	Cvc string `json:"cvc"`
}
