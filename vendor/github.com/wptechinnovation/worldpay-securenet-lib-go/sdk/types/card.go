package types

// Card represents a payment card
type Card struct {
	Number         string   `json:"number"`
	ExpirationDate string   `json:"expirationDate"`
	Address        *Address `json:"address"`
}
