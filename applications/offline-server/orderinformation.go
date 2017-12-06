package main

import (
	"time"
)

// OrderRequest A request to create an order
type OrderInformation struct {
	MaskedCard string `json:"maskedCard"`

	TotalPrice int `json:"amount"`

	CurrencyCode string `json:"currencyCode"`

	OrderDescription string `json:"orderDescription"`

	TransactionDateAndTime time.Time `json:"transactionDateAndTime"`
}
