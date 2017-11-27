package core

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/configuration"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/hte"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp/securenet"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/servicediscovery"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types/event"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/utils"
)

const (

	// BroadcastStepSleep The amount of time to sleep between sending each broadcast message (Milliseconds)
	BroadcastStepSleep = 5000
	// BroadcastPort The port to broadcast messages on
	BroadcastPort = 8980
	// HteSvcURLPrefix HTE REST API Url prefix - can be empty
	HteSvcURLPrefix = ""
	// UUIDFilePath Path to store devie UUID once created
	UUIDFilePath = "uuid.txt"
	// HteSvcPort Port that the HTE REST API listens on
	HteSvcPort = 64521
	// WPOnlineAPIEndpoint Worldpay online API endpoint
	WPOnlineAPIEndpoint = "https://api.worldpay.com/v1"
	// HteClientScheme HTE REST API Scheme typically http:// or https://
	HteClientScheme = "http://"
)

// SDKFactory for creating WPWithin instances. // TODO Needs to be reworked so can be partially implemented.
type SDKFactory interface {
	GetDevice(name, description, interfaceAddr string, cfg *configuration.WPWithin) (*types.Device, error)
	GetPSPMerchant(pspConfig map[string]string) (psp.PSP, error)
	GetPSPClient(pspConfig map[string]string) (psp.PSP, error)
	GetSvcBroadcaster(ipv4Address string) (servicediscovery.Broadcaster, error)
	GetSvcScanner() (servicediscovery.Scanner, error)
	GetHTE(device *types.Device, psp psp.PSP, ipv4Address string, scheme string, hteCredential *hte.Credential, om hte.OrderManager, hteSvcHandler *hte.ServiceHandler) (hte.Service, error)
	GetOrderManager() (hte.OrderManager, error)
	GetHTEClient() (hte.Client, error)
	GetHTEClientHTTP() (hte.ClientHTTP, error)
	GetHTEServiceHandler(device *types.Device, psp psp.PSP, credential *hte.Credential, orderManager hte.OrderManager, eventHandler event.Handler) *hte.ServiceHandler
}

// SDKFactoryImpl implementation of SDKFactory
type SDKFactoryImpl struct{}

// NewSDKFactory create a new SDKFactory
func NewSDKFactory() (SDKFactory, error) {

	return &SDKFactoryImpl{}, nil
}

// GetDevice create a device with Name and Description
func (factory *SDKFactoryImpl) GetDevice(name, description, interfaceAddr string, cfg *configuration.WPWithin) (*types.Device, error) {

	var deviceGUID string

	if b, _ := utils.FileExists(UUIDFilePath); b {

		_deviceGUID, err := utils.ReadLocalUUID(UUIDFilePath)

		if err != nil {

			return nil, fmt.Errorf("Could not read UUID file (%s). Try deleting it. %q", UUIDFilePath, err.Error())
		}

		deviceGUID = _deviceGUID

	} else {

		_deviceGUID, err := utils.NewUUID()

		if err != nil {

			return nil, fmt.Errorf("Unable to create new UUID: %q", err.Error())
		}

		deviceGUID = _deviceGUID

		if err := utils.WriteString(UUIDFilePath, deviceGUID, true); err != nil {

			return nil, fmt.Errorf("Could not save UUID to %s", UUIDFilePath)
		}
	}

	var deviceAddress net.IP
	deviceAddress = net.ParseIP(interfaceAddr)
	if deviceAddress == nil {
		deviceAddress, _ = utils.FirstExternalIPv4()
	}
	d, e := types.NewDevice(name, description, deviceGUID, deviceAddress.String(), "GBP")

	return d, e
}

// GetPSPMerchant get a new PSP implementation in context of Merchant i.e. client/service keys are set
func (factory *SDKFactoryImpl) GetPSPMerchant(pspConfig map[string]string) (psp.PSP, error) {

	if pspConfig == nil {

		return nil, errors.New("PSP Config map must be set")
	}

	switch pspConfig[psp.CfgPSPName] {

	case onlineworldpay.PSPName:
		return onlineworldpay.NewMerchant(pspConfig[onlineworldpay.CfgMerchantClientKey], pspConfig[onlineworldpay.CfgMerchantServiceKey], pspConfig[onlineworldpay.CfgAPIEndpoint])

	case securenet.PSPName:
		devID, err := strconv.Atoi(pspConfig[securenet.CfgDeveloperID])

		if err != nil {

			return nil, fmt.Errorf("Error parsing developerID as int")
		}

		return securenet.NewSecureNetMerchant(pspConfig[securenet.CfgSecureNetID], pspConfig[securenet.CfgSecureKey], pspConfig[securenet.CfgPublicKey], pspConfig[securenet.CfgAPIEndpoint], pspConfig[securenet.CfgAppVersion], int32(devID), pspConfig[securenet.CfgHTTPProxy])
	}

	return nil, fmt.Errorf("Unknown PSP: %s", pspConfig[psp.CfgPSPName])
}

// GetPSPClient get a new PSP implementation in context of a client i.e. only the endpoint is set
func (factory *SDKFactoryImpl) GetPSPClient(pspConfig map[string]string) (psp.PSP, error) {

	if pspConfig == nil {

		return nil, errors.New("PSP Config map must be set")
	}

	switch pspConfig[psp.CfgPSPName] {

	case "worldpayonlinepayments":
		return onlineworldpay.NewClient(pspConfig[onlineworldpay.CfgAPIEndpoint])

	case "securenet":
		devID, err := strconv.Atoi(pspConfig[securenet.CfgDeveloperID])

		if err != nil {

			return nil, fmt.Errorf("Error parsing developerID as int")
		}

		return securenet.NewSecureNetConsumer(pspConfig[securenet.CfgAPIEndpoint], pspConfig[securenet.CfgAppVersion], int32(devID), pspConfig[securenet.CfgHTTPProxy])
	}

	return nil, fmt.Errorf("Unknown PSP: %s", pspConfig[psp.CfgPSPName])
}

// GetSvcBroadcaster get an instance of service broadcaster
func (factory *SDKFactoryImpl) GetSvcBroadcaster(ipv4Address string) (servicediscovery.Broadcaster, error) {

	return servicediscovery.NewBroadcaster(ipv4Address, BroadcastPort, BroadcastStepSleep)
}

// GetSvcScanner get an instance of service scanner
func (factory *SDKFactoryImpl) GetSvcScanner() (servicediscovery.Scanner, error) {

	return servicediscovery.NewScanner(BroadcastPort, BroadcastStepSleep)
}

// GetHTE get an instance of HTE
func (factory *SDKFactoryImpl) GetHTE(device *types.Device, psp psp.PSP, ipv4Address string, scheme string, hteCredential *hte.Credential, om hte.OrderManager, hteSvcHandler *hte.ServiceHandler) (hte.Service, error) {

	return hte.NewService(device, psp, ipv4Address, HteSvcURLPrefix, scheme, HteSvcPort, hteCredential, om, hteSvcHandler)
}

// GetOrderManager get an instance of OrderManager
func (factory *SDKFactoryImpl) GetOrderManager() (hte.OrderManager, error) {

	return hte.NewOrderManager()
}

// GetHTEClient get an instance of HTEClient
func (factory *SDKFactoryImpl) GetHTEClient() (hte.Client, error) {

	return nil, nil
}

// GetHTEClientHTTP get an instance of HTEClientHTTP
func (factory *SDKFactoryImpl) GetHTEClientHTTP() (hte.ClientHTTP, error) {

	return hte.NewHTEClientHTTP()
}

// GetHTEServiceHandler get an instance of HTE Service Handler
func (factory *SDKFactoryImpl) GetHTEServiceHandler(device *types.Device, psp psp.PSP, credential *hte.Credential, orderManager hte.OrderManager, eventHandler event.Handler) *hte.ServiceHandler {

	return hte.NewServiceHandler(device, psp, credential, orderManager, eventHandler)
}
