package service

// Connection defines functions to enable connectivity with the api endpoints
type Connection interface {

	// Post a message to the worldpay server
	// body is to be marshalled into JSON and submitted as the POST body
	// url is appended to the connection APIEndpoint
	// auth specifies whether the authorisation header (service key) should be added (true) or not (false)
	Post(body []byte, url string, auth bool) ([]byte, error)
}
