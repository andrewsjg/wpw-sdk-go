package servicediscovery

import "github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"

// Scanner defines function for discovering devices on a network
type Scanner interface {
	ScanForServices(timeout int) (map[string]types.BroadcastMessage, error)
	ScanForService(timeout int, serviceName string) (*types.BroadcastMessage, error)
	StopScanner()
}
