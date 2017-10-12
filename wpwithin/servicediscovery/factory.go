package servicediscovery

// NewScanner creates a new instance Scanner
func NewScanner(port, stepSleep int) (Scanner, error) {

	result := &scannerImpl{
		run:       false,
		stepSleep: stepSleep,
	}

	comm, err := NewUDPComm(port, "udp4")

	if err != nil {

		return nil, err
	}

	result.comm = comm

	return result, nil
}

// NewBroadcaster create a new instance of Broadcaster
func NewBroadcaster(bcastaddr string, port int, stepSleep int) (Broadcaster, error) {

	result := &broadcasterImpl{

		stepSleep: stepSleep,
		run:       false,
		host:      bcastaddr, //net.IPv4bcast.To4().String(),
		port:      port,
	}

	comm, err := NewUDPComm(port, "udp4")

	if err != nil {

		return nil, err
	}

	result.comm = comm

	return result, nil
}
