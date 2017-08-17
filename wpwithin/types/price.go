package types

// Price represents price
type Price struct {
	ID              int           `json:"priceID"`
	Description     string        `json:"priceDescription"`
	PricePerUnit    *PricePerUnit `json:"pricePerUnit"`
	UnitID          int           `json:"unitID"`
	UnitDescription string        `json:"unitDescription"`
}

// NewPrice create a new instance of Price
func NewPrice() (*Price, error) {

	result := &Price{}

	return result, nil
}
