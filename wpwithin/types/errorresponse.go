package types

// ErrorResponse error response message
type ErrorResponse struct {
	HTTPStatusCode int    `json:"HTTPStatusCode"`
	Message        string `json:"message"`
	ErrorCode      int    `json:"errorCode"`
}
