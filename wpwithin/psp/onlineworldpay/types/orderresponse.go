package types

// OrderResponse details of the response to creating an order
type OrderResponse struct {
	OrderCode string `json:"orderCode"`

	Token string `json:"token"`

	OrderDescription string `json:"orderDescription"`

	Amount int32 `json:"amount"`

	CurrencyCode string `json:"currencyCode"`

	PaymentStatus string `json:"paymentStatus"`

	PaymentResponse OrderResponsePaymentResponse `json:"paymentResponse"`

	CustomerOrderCode string `json:"customerOrderCode"`

	Environment string `json:"environment"`

	RiskScore OrderResponseRiskScore `json:"riskScore"`
}
