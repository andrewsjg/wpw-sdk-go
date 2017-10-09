package cardnotpresent

import "time"

// ChargeTokenResponse is a response to a request to charge a token when card is not present
type ChargeTokenResponse struct {
	Result       string                 `json:"result"`
	ResponseCode int                    `json:"responseCode"`
	Messages     string                 `json:"message"`
	Transaction  ChargeTokenTransaction `json:"transaction"`
}

// ChargeTokenTransaction is the transaction object returned from the ChargeToken call
type ChargeTokenTransaction struct {
	SecureNetID           int     `json:"secureNetId"`
	TransactionType       string  `json:"transactionType"`
	CustomerID            string  `json:"customerId"`
	OrderID               string  `json:"orderId"`
	TransactionID         int     `json:"transactionId"`
	AuthorizationCode     string  `json:"authorizationCode"`
	AuthorizedAmount      float64 `json:"authorizedAmount"`
	AllowedPartialCharges bool    `json:"allowedPartialCharges"`
	PaymentTypeCode       string  `json:"paymentTypeCode"`
	PaymentTypeResult     string  `json:"paymentTypeResult"`
	Level2Valid           bool    `json:"level2Valid"`
	Level3Valid           bool    `json:"level3Valid"`
	TransactionData       struct {
		Date   time.Time `json:"date"`
		Amount float64   `json:"amount"`
	} `json:"transactionData"`
	SettlementData      interface{} `json:"settlementData"`
	VaultData           interface{} `json:"vaultData"`
	CreditCardType      string      `json:"creditCardType"`
	CardNumber          string      `json:"cardNumber"`
	AvsCode             string      `json:"avsCode"`
	AvsResult           string      `json:"avsResult"`
	CardHolderFirstName string      `json:"cardHolder_FirstName"`
	CardHolderLastName  string      `json:"cardHolder_LastName"`
	ExpirationDate      string      `json:"expirationDate"`
	BillAddress         struct {
		Line1     string      `json:"line1"`
		City      string      `json:"city"`
		State     string      `json:"state"`
		Zip       string      `json:"zip"`
		Country   string      `json:"country"`
		Company   string      `json:"company"`
		Phone     string      `json:"phone"`
		FirstName interface{} `json:"firstName"`
		LastName  interface{} `json:"lastName"`
	} `json:"billAddress"`
	Email                string      `json:"email"`
	EmailReceipt         bool        `json:"emailReceipt"`
	CardCodeCode         string      `json:"cardCodeCode"`
	CardCodeResult       string      `json:"cardCodeResult"`
	AccountName          interface{} `json:"accountName"`
	AccountType          interface{} `json:"accountType"`
	AccountNumber        interface{} `json:"accountNumber"`
	CheckNumber          interface{} `json:"checkNumber"`
	TraceNumber          interface{} `json:"traceNumber"`
	SurchargeAmount      float64     `json:"surchargeAmount"`
	CashbackAmount       float64     `json:"cashbackAmount"`
	FnsNumber            interface{} `json:"fnsNumber"`
	VoucherNumber        interface{} `json:"voucherNumber"`
	FleetCardInfo        interface{} `json:"fleetCardInfo"`
	Gratuity             float64     `json:"gratuity"`
	IndustrySpecificData string      `json:"industrySpecificData"`
	MarketSpecificData   string      `json:"marketSpecificData"`
	NetworkCode          string      `json:"networkCode"`
	AdditionalAmount     float64     `json:"additionalAmount"`
	AdditionalData1      interface{} `json:"additionalData1"`
	AdditionalData2      interface{} `json:"additionalData2"`
	AdditionalData3      interface{} `json:"additionalData3"`
	AdditionalData4      string      `json:"additionalData4"`
	AdditionalData5      string      `json:"additionalData5"`
	Method               string      `json:"method"`
	ResponseText         string      `json:"responseText"`
	ImageResult          interface{} `json:"imageResult"`
}
