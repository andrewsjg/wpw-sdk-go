package tokenization

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/service"
)

// Service defines Tokenization related functions
type Service interface {
	TokenizeCard(card *TokenizeCardRequest) (*TokenizeCardResponse, error)
}

type serviceImpl struct {
	Connection service.Connection
}

// NewService returns a new implemenetation of Tokenization service
func NewService(connection service.Connection) (Service, error) {

	result := serviceImpl{}
	result.Connection = connection

	return &result, nil
}

func (svc *serviceImpl) TokenizeCard(request *TokenizeCardRequest) (*TokenizeCardResponse, error) {

	log.Debug("Begin TokenizeCard()")

	if request == nil {

		return nil, errors.New("Request cannot be nil")
	}

	bodyBytes, err := json.Marshal(request)

	if err != nil {

		return nil, err
	}

	respBytes, err := svc.Connection.Post(bodyBytes, "/PreVault/Card", true)

	if err != nil {

		return nil, err
	}

	var respObject *TokenizeCardResponse

	err = json.Unmarshal(respBytes, &respObject)

	log.Debug("End TokenizeCard()")

	return respObject, err
}
