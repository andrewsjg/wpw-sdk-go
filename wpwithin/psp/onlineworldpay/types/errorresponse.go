package types

// ErrorResponse Details of a response from the API that is an error
type ErrorResponse struct {
	HTTPStatusCode  int    `json:"HttpStatusCode"`
	CustomCode      string `json:"customCode"`
	Message         string `json:"message"`
	Description     string `json:"description"`
	ErrorHelpURL    string `json:"ErrorHelpUrl"`
	OriginalRequest string `json:"originalRequest"`
}
