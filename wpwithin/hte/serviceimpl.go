package hte

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
)

// ServiceImpl Concrete Implementation of HTE service
type ServiceImpl struct {
	router        *mux.Router
	_IPv4Address  string
	_Port         int
	_URLPrefix    string
	_Scheme       string
	routes        []Route
	handler       *ServiceHandler
	HTECredential *Credential
	orderManager  OrderManager
}

// Route a route to a function
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// NewService crate a new instance of Service
func NewService(device *types.Device, psp psp.PSP, ip, prefix, scheme string, port int, hteCredential *Credential, orderManager OrderManager, svcHandler *ServiceHandler) (Service, error) {

	service := &ServiceImpl{}

	service.handler = svcHandler

	service._URLPrefix = prefix
	service._Port = port
	service._IPv4Address = ip
	service._Scheme = scheme
	service.orderManager = orderManager

	service.setupRoutes()

	service.router = mux.NewRouter().StrictSlash(true)

	for _, route := range service.routes {

		service.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return service, nil
}

// Start the service
func (service *ServiceImpl) Start() error {

	return http.ListenAndServe(fmt.Sprintf(":%d", service.Port()), service.router)
}

// setupRoutes configures the services routes
func (service *ServiceImpl) setupRoutes() {

	service.routes = []Route{
		Route{
			"Service Discovery",
			"GET",
			fmt.Sprintf("%s/service/discover", service.URLPrefix()),
			service.handler.ServiceDiscovery,
		},
		Route{
			"Service Price Request",
			"GET",
			fmt.Sprintf("%s/service/{service_id}/prices", service.URLPrefix()),
			service.handler.ServicePrices,
		},
		Route{
			"Service Total Price Request",
			"POST",
			fmt.Sprintf("%s/service/{service_id}/requestTotal", service.URLPrefix()),
			service.handler.ServiceTotalPrice,
		},
		Route{
			"Payment",
			"POST",
			fmt.Sprintf("%s/payment", service.URLPrefix()),
			service.handler.Payment,
		},
		Route{
			"Service Delivery Begin",
			"POST",
			fmt.Sprintf("%s/service/{service_id}/delivery/begin", service.URLPrefix()),
			service.handler.ServiceDeliveryBegin,
		},
		Route{
			"Service Delivery End",
			"POST",
			fmt.Sprintf("%s/service/{service_id}/delivery/end", service.URLPrefix()),
			service.handler.ServiceDeliveryEnd,
		},
	}
}

// IPAddr get the service IP Address
func (service *ServiceImpl) IPAddr() string {

	return service._IPv4Address
}

// Port get the service port
func (service *ServiceImpl) Port() int {

	return service._Port
}

// URLPrefix get the service URL prefix
func (service *ServiceImpl) URLPrefix() string {

	return service._URLPrefix
}

// Scheme get the service scheme
func (service *ServiceImpl) Scheme() string {

	return service._Scheme
}
