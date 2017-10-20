package types

// PaymentVaultToken represents a tokenised payment credential that is stored in secure vault
type PaymentVaultToken struct {
	CustomerID      string      `json:"customerId"`
	PaymentMethodID string      `json:"paymentMethodId"`
	PublicKey       string      `json:"publicKey"`
	PaymentType     PaymentType `json:"paymentType"`
}

// PaymentType defines type of payment being made
type PaymentType string

const (
	// Check payment type
	Check PaymentType = "CHECK"
	// DebitCard debit card
	DebitCard PaymentType = "DEBIT_CARD"
	// CreditCard credit card
	CreditCard PaymentType = "CREDIT_CARD"
	// FleetCard fleet card
	FleetCard PaymentType = "FLEET_CARD"
	// StoredValue stored value
	StoredValue PaymentType = "STORED_VALUE"
	// Unknown unknown
	Unknown PaymentType = "UNKNOWN"
)
