package hte

// Service HTE Service
type Service interface {

	// Start the HTE service to allow consumer connect
	Start() error
	// Setup the routes i.e. The urls that map to functions
	setupRoutes()
	// Listening IP address
	IPAddr() string
	// Listening port
	Port() int
	// Url prefix of route URLs
	URLPrefix() string
	// Scheme e.g. HTTP/HTTPS
	Scheme() string
}
