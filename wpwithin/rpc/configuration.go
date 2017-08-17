package rpc

// Configuration RPC configuration
type Configuration struct {
	Protocol     string
	Framed       bool
	Buffered     bool
	Host         string
	Port         int
	Secure       bool
	BufferSize   int
	Loglevel     string
	Logfile      string
	CallbackPort int
}
