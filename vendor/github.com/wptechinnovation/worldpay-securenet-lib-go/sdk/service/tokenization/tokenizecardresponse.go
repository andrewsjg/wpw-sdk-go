package tokenization

// TokenizeCardResponse Contains information returned from a request to tokenize a payment card
type TokenizeCardResponse struct {
	HTTPStatusCode int
	Result         string `json:"result"`
	ResponseCode   int32  `json:"responseCode"`
	CustomerID     string `json:"customerId"`
	Token          string `json:"token"`
}
