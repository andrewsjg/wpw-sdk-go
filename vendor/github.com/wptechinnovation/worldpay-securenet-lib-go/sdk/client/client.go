package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/wptechinnovation/worldpay-securenet-lib-go/sdk/service/cardnotpresent"
	"github.com/wptechinnovation/worldpay-securenet-lib-go/sdk/service/tokenization"
)

// Client enables interaction with the Worldpay API
type Client interface {
	CardNotPresentService() cardnotpresent.Service
	TokenizationService() tokenization.Service
}

type clientImpl struct {
	_cardNotPresentService cardnotpresent.Service
	_tokenizationService   tokenization.Service
}

// New returns a new implemention of SDK Client
func New(apiEndpoint, appVersion, secureNetID, secureKey string, proxy string) (Client, error) {

	log.WithFields(log.Fields{"apiEndpoint": apiEndpoint, "appVersion": appVersion, "secureNetID": secureNetID, "secureKey": secureKey}).Debug("begin client.New()")

	connection, err := NewConnection(apiEndpoint, appVersion, secureNetID, secureKey, 60, true, proxy)

	if err != nil {

		return nil, err
	}

	log.Debug("Did create new connection")

	tokenSvc, err := tokenization.NewService(connection)

	if err != nil {

		log.WithField("Error", err.Error()).Error("Error creating new tokenization service")

		return nil, err
	}

	log.Debug("Did create new tokenization service")

	cnpSvc, err := cardnotpresent.NewService(connection)

	if err != nil {

		log.WithField("Error", err.Error()).Error("Error creating new card not present service")

		return nil, err
	}

	log.Debug("Did create card not present service")

	result := &clientImpl{}
	result._tokenizationService = tokenSvc
	result._cardNotPresentService = cnpSvc

	defer log.Debug("end client.New()")

	return result, nil
}

func (client *clientImpl) CardNotPresentService() cardnotpresent.Service {

	log.Debug("clientImpl.CardNotPresentService()")

	return client._cardNotPresentService
}

func (client *clientImpl) TokenizationService() tokenization.Service {

	log.Debug("clientImpl.TokenizationService()")

	return client._tokenizationService
}
