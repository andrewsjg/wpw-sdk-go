package types

// Address related to a card
type Address struct {
	Line1   string `json:"line1"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
	Company string `json:"company"`
	Phone   string `json:"phone"`
}
