package servicediscovery

import "github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"

// Broadcaster defines functionality to broadcast a devices presence on a network
type Broadcaster interface {
	StartBroadcast(msg types.BroadcastMessage, timeoutMillis int) error
	StopBroadcast() error
}
