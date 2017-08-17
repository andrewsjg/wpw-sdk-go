package rpc

import (
	"crypto/tls"
	"fmt"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
	log "github.com/Sirupsen/logrus"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/rpc/wpthrift/gen-go/wpthrift"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/types/event"
)

// ServiceImpl impementation of RPC service
type ServiceImpl struct {
	wpWithin         wpwithin.WPWithin
	transportFactory thrift.TTransportFactory
	protocolFactory  thrift.TProtocolFactory
	host             string
	port             int
	secure           bool
	callbackClient   event.Handler
}

// NewService create a new instance of Service
func NewService(config Configuration, wpWithin wpwithin.WPWithin) (Service, error) {

	log.WithField("Config", fmt.Sprintf("%+v", config)).Debug("begin rpc.ServiceImpl.NewService()")

	defer log.Debug("end rpc.ServiceImpl.NewService()")

	// Validate configuration host - should be local only as we don't want to receive commands
	// from remote hosts
	if !strings.EqualFold(config.Host, "127.0.0.1") && !strings.EqualFold(config.Host, "localhost") {

		return nil, fmt.Errorf("Invalid configuration.Host provided (%s) - this service will only listen on local interface i.e. 127.0.0.1 or localhost", config.Host)
	}

	var protocolFactory thrift.TProtocolFactory
	switch config.Protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		return nil, fmt.Errorf("Invalid protocol specified: %s\n", config.Protocol)
	}

	var transportFactory thrift.TTransportFactory
	if config.Buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(config.BufferSize)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if config.Framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	result := new(ServiceImpl)
	result.transportFactory = transportFactory
	result.protocolFactory = protocolFactory
	result.host = config.Host
	result.port = config.Port
	result.secure = config.Secure
	result.wpWithin = wpWithin
	// Only setup callbacks if there is a callback port specified
	if config.CallbackPort > 0 {

		log.Debug("CallbackPort > 0, Setting up RPC Callback")
		log.Debug("Will attempt to create new EventSender")

		cb, err := NewEventSender(config)

		if err != nil {

			return nil, err
		}

		log.Debug("Did create new EventSender.")
		result.callbackClient = cb
	}

	return result, nil
}

// Start the RPC service
func (svc *ServiceImpl) Start() error {

	log.Debug("begin rpc.ServiceImpl.start()")

	defer log.Debug("end rpc.ServiceImpl.start()")

	strAddr := fmt.Sprintf("%s:%d", svc.host, svc.port)

	var transport thrift.TServerTransport
	var err error
	if svc.secure {
		cfg := new(tls.Config)
		if cert, _err := tls.LoadX509KeyPair("server.crt", "server.key"); _err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return _err
		}
		transport, err = thrift.NewTSSLServerSocket(strAddr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(strAddr)
	}

	if err != nil {
		return err
	}

	handler := NewWPWithinHandler(svc.wpWithin, svc.callbackClient)
	processor := wpthrift.NewWPWithinProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, svc.transportFactory, svc.protocolFactory)

	log.Debugf("Starting the rpc server on...: %s", strAddr)

	return server.Serve()
}
