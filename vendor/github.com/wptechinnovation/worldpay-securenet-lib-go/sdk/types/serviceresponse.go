package types

// ServiceResponse Contains information returned from a request to Worldpay API
type ServiceResponse struct {
	HTTPStatusCode int
	Result         string      `json:"result"`
	ResponseCode   int32       `json:"responseCode"`
	Message        string      `json:"message"`
	Transaction    interface{} `json:"transaction"`
}
