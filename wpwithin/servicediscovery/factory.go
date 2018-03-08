package servicediscovery

import (
	"net"

	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/wpwerrors"
)

// NewScanner creates a new instance Scanner
func NewScanner(port, stepSleep int) (Scanner, error) {

	result := &scannerImpl{
		run:       false,
		stepSleep: stepSleep,
	}

	comm, err := NewUDPComm(port, "udp4")

	if err != nil {

		return nil, wpwerrors.GetError(wpwerrors.NEWUDPCOMMERR)
	}

	result.comm = comm

	return result, nil
}

// NewBroadcaster create a new instance of Broadcaster
func NewBroadcaster(bcastaddr string, port int, stepSleep int) (Broadcaster, error) {
	// TODO: bcastaddr is not used

	result := &broadcasterImpl{

		stepSleep:     stepSleep,
		run:           false,
		interfaceAddr: bcastaddr,
		host:          net.IPv4bcast.To4().String(),
		port:          port,
	}

	comm, err := NewUDPComm(port, "udp4")

	if err != nil {

		return nil, wpwerrors.GetError(wpwerrors.NEWUDPCOMMERR)
	}

	result.comm = comm

	return result, nil
}
