package types

// PricePerUnit ...
type PricePerUnit struct {
	Amount       int    `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}
