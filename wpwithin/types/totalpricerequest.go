package types

// TotalPriceRequest HTTP Message
type TotalPriceRequest struct {
	ClientID              string `json:"clientID"`
	SelectedNumberOfUnits int    `json:"selectedNumberOfUnits"`
	SelectedPriceID       int    `json:"selectedPriceID"`
}
