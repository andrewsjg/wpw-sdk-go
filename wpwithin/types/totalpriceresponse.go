package types

// TotalPriceResponse Details of a the price response from HTE server. i.e. Total amount to
// be paid for selected service
type TotalPriceResponse struct {
	ServerID           string `json:"serverID"`
	ClientID           string `json:"clientID"`
	PriceID            int    `json:"priceID"`
	UnitsToSupply      int    `json:"unitsToSupply"`
	TotalPrice         int    `json:"totalPrice"`
	CurrencyCode       string `json:"currencyCode"`
	PaymentReferenceID string `json:"paymentReferenceID"`
	MerchantClientKey  string `json:"merchantClientKey"`
}
