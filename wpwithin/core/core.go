package core

import (
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/configuration"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/hte"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/psp"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/servicediscovery"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/types"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/types/event"
)

// Core This acts as a container for dependencies of the SDK
type Core struct {
	Device         *types.Device
	Psp            psp.PSP
	SvcBroadcaster servicediscovery.Broadcaster
	SvcScanner     servicediscovery.Scanner
	HTE            hte.Service
	HCECard        *types.HCECard
	OrderManager   hte.OrderManager
	HTEClient      hte.Client
	EventHandler   event.Handler
	Configuration  configuration.WPWithin
}

// NewCore Create a new Core
func NewCore() (*Core, error) {

	result := &Core{}

	return result, nil
}

// SetDevice set the device
func (core *Core) SetDevice(device *types.Device) {

	core.Device = device
}

// SetPsp set the PSP
func (core *Core) SetPsp(psp psp.PSP) {

	core.Psp = psp
}

// SetSvcBroadcaster set the service broadcaster
func (core *Core) SetSvcBroadcaster(svcBroadcaster servicediscovery.Broadcaster) {

	core.SvcBroadcaster = svcBroadcaster
}

// SetSvcScanner set the service scanner
func (core *Core) SetSvcScanner(serviceScanner servicediscovery.Scanner) {

	core.SvcScanner = serviceScanner
}

// SetHTE set the HTE Service
func (core *Core) SetHTE(hteService hte.Service) {

	core.HTE = hteService
}

// SetHCECard set the HCE card
func (core *Core) SetHCECard(hceCard *types.HCECard) {

	core.HCECard = hceCard
}

// SetOrderManager set the order manager
func (core *Core) SetOrderManager(orderManager hte.OrderManager) {

	core.OrderManager = orderManager
}

// SetHTEClient set the HTE client
func (core *Core) SetHTEClient(hteClient hte.Client) {

	core.HTEClient = hteClient
}
