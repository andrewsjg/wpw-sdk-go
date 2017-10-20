package cardnotpresent

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/service"
)

// Service defines functions related card not present transactions
type Service interface {
	ChargeUsingToken(request *ChargeTokenRequest) (*ChargeTokenResponse, error)
}

type serviceImpl struct {
	connection service.Connection
}

// NewService returns a new implemenetation of Card Not Present service
func NewService(connection service.Connection) (Service, error) {

	result := serviceImpl{}
	result.connection = connection

	return &result, nil
}

func (svc *serviceImpl) ChargeUsingToken(request *ChargeTokenRequest) (*ChargeTokenResponse, error) {

	logrus.Debug("Begin ChargeUsingToken()")

	if request == nil {

		return nil, errors.New("Request cannot be nil")
	}

	bodyBytes, err := json.Marshal(request)

	if err != nil {

		return nil, err
	}

	respBytes, err := svc.connection.Post(bodyBytes, "/Payments/Charge", true)

	if err != nil {

		return nil, err
	}

	var respObject *ChargeTokenResponse

	err = json.Unmarshal(respBytes, &respObject)

	logrus.Debug("End ChargeUsingToken()")

	return respObject, err
}
