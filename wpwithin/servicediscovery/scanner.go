package servicediscovery

import "github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/types"

// Scanner defines function for discovering devices on a network
type Scanner interface {
	ScanForServices(timeout int) (map[string]types.BroadcastMessage, error)

	StopScanner()
}
