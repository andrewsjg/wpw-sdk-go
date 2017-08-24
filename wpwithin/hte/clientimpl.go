package hte

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"

	log "github.com/Sirupsen/logrus"
)

// Concrete implementation of HTEClient
type clientImpl struct {
	scheme     string
	host       string
	port       int
	urlPrefix  string
	clientID   string
	baseURL    string
	httpClient ClientHTTP
}

// NewClient Initialise a new HTE Client
func NewClient(scheme, host string, port int, urlPrefix string, clientID string, httpClient ClientHTTP) (Client, error) {

	if host == "" {

		return nil, errors.New("host cannot be empty")
	}

	if port < PortRangeMin || port > PortRangeMax {

		return nil, fmt.Errorf("Port number cannot exceed range [%d - %d]", PortRangeMin, PortRangeMax)
	}

	// Do not need to validate against empty urlPrefix as it can actually be empty

	if clientID == "" {

		return nil, errors.New("clientId cannot be empty")
	}

	result := &clientImpl{}
	result.host = host
	result.port = port
	result.urlPrefix = urlPrefix
	result.clientID = clientID
	result.httpClient = httpClient

	// Compose base url and ensure there are no duplicate slashes from passed in urlPrefix
	result.baseURL = fmt.Sprintf("%s%s:%d%s", scheme, host, port, urlPrefix)

	return result, nil
}

func (client *clientImpl) DiscoverServices() (types.ServiceListResponse, error) {

	url := fmt.Sprintf("%s/service/discover", client.baseURL)

	response, err := client.httpClient.Get(url)

	if err != nil {

		return types.ServiceListResponse{}, err
	}

	svcDetails := types.ServiceListResponse{}

	err = json.Unmarshal(response, &svcDetails)

	if err != nil {

		return types.ServiceListResponse{}, err
	}

	return svcDetails, nil
}

func (client *clientImpl) GetPrices(serviceID int) (types.ServicePriceResponse, error) {

	url := fmt.Sprintf("%s/service/%d/prices", client.baseURL, serviceID)

	response, err := client.httpClient.Get(url)

	if err != nil {

		return types.ServicePriceResponse{}, err
	}

	svcPriceResponse := types.ServicePriceResponse{}

	err = json.Unmarshal(response, &svcPriceResponse)

	if err != nil {

		return types.ServicePriceResponse{}, err
	}

	return svcPriceResponse, nil
}

func (client *clientImpl) NegotiatePrice(serviceID, priceID, numberOfUnits int) (types.TotalPriceResponse, error) {

	url := fmt.Sprintf("%s/service/%d/requestTotal", client.baseURL, serviceID)

	req := types.TotalPriceRequest{
		ClientID:              client.clientID,
		SelectedNumberOfUnits: numberOfUnits,
		SelectedPriceID:       priceID,
	}

	jsonReq, err := json.Marshal(req)

	if err != nil {

		return types.TotalPriceResponse{}, err
	}

	bytesResp, httpStatus, err := client.httpClient.PostJSON(url, jsonReq)

	if err != nil {

		return types.TotalPriceResponse{}, err
	} else if httpStatus != http.StatusOK {

		errorResp := types.ErrorResponse{}

		err = json.Unmarshal(bytesResp, &errorResp)

		if err != nil {

			return types.TotalPriceResponse{}, err
		}

		return types.TotalPriceResponse{}, fmt.Errorf("%d - %s (%d)", errorResp.HTTPStatusCode, errorResp.Message, errorResp.ErrorCode)

	} else {

		priceResp := types.TotalPriceResponse{}

		err = json.Unmarshal(bytesResp, &priceResp)

		if err != nil {

			return types.TotalPriceResponse{}, err
		}

		return priceResp, nil
	}
}

func (client *clientImpl) MakeHtePayment(paymentReferenceID, clientID, clientToken string) (types.PaymentResponse, error) {

	url := fmt.Sprintf("%s/payment", client.baseURL)

	requestBody := types.PaymentRequest{
		ClientID:           clientID,
		ClientToken:        clientToken,
		PaymentReferenceID: paymentReferenceID,
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {

		return types.PaymentResponse{}, err
	}

	byteResp, httpStatus, err := client.httpClient.PostJSON(url, jsonBody)

	if err != nil {

		return types.PaymentResponse{}, err
	} else if httpStatus != http.StatusOK {

		errorResponse := types.ErrorResponse{}

		err = json.Unmarshal(byteResp, &errorResponse)

		if err != nil {

			return types.PaymentResponse{}, err
		}

		return types.PaymentResponse{}, fmt.Errorf("%d - %s (%d)", errorResponse.HTTPStatusCode, errorResponse.Message, errorResponse.ErrorCode)

	} else {

		paymentResponse := types.PaymentResponse{}

		err = json.Unmarshal(byteResp, &paymentResponse)

		if err != nil {

			return types.PaymentResponse{}, err
		}

		return paymentResponse, nil
	}
}

func (client *clientImpl) StartDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsToSupply int) (types.BeginServiceDeliveryResponse, error) {

	log.WithFields(log.Fields{"serviceId": serviceID, "serviceDeliveryToken": serviceDeliveryToken, "unitsToSupply": unitsToSupply}).Debug("hte.clientImpl.StartDelivery()")

	defer log.Debug("end hte.clientImpl.EndDelivery()")

	url := fmt.Sprintf("%s/service/%d/delivery/begin", client.baseURL, serviceID)

	requestBody := types.BeginServiceDeliveryRequest{
		ClientID:             client.clientID,
		ServiceDeliveryToken: serviceDeliveryToken,
		UnitsToSupply:        unitsToSupply,
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {

		return types.BeginServiceDeliveryResponse{}, err
	}

	byteResp, httpStatus, err := client.httpClient.PostJSON(url, jsonBody)

	if err != nil {
		log.Errorf("Error posting JSON request. Error = %s", err.Error())
		return types.BeginServiceDeliveryResponse{}, err
	} else if httpStatus != http.StatusOK {

		errorResponse := types.ErrorResponse{}

		err = json.Unmarshal(byteResp, &errorResponse)

		if err != nil {

			return types.BeginServiceDeliveryResponse{}, err
		}

		return types.BeginServiceDeliveryResponse{}, fmt.Errorf("%d - %s (%d)", errorResponse.HTTPStatusCode, errorResponse.Message, errorResponse.ErrorCode)

	} else {

		deliveryResponse := types.BeginServiceDeliveryResponse{}

		err = json.Unmarshal(byteResp, &deliveryResponse)

		if err != nil {

			return types.BeginServiceDeliveryResponse{}, err
		}

		return deliveryResponse, nil
	}
}

func (client *clientImpl) EndDelivery(serviceID int, serviceDeliveryToken types.ServiceDeliveryToken, unitsReceived int) (types.EndServiceDeliveryResponse, error) {

	url := fmt.Sprintf("%s/service/%d/delivery/end", client.baseURL, serviceID)

	requestBody := types.EndServiceDeliveryRequest{
		ClientID:             client.clientID,
		ServiceDeliveryToken: serviceDeliveryToken,
		UnitsReceived:        unitsReceived,
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {

		return types.EndServiceDeliveryResponse{}, err
	}

	byteResp, httpStatus, err := client.httpClient.PostJSON(url, jsonBody)

	if err != nil {

		return types.EndServiceDeliveryResponse{}, err
	} else if httpStatus != http.StatusOK {

		errorResponse := types.ErrorResponse{}

		err = json.Unmarshal(byteResp, &errorResponse)

		if err != nil {

			return types.EndServiceDeliveryResponse{}, err
		}

		return types.EndServiceDeliveryResponse{}, fmt.Errorf("%d - %s (%d)", errorResponse.HTTPStatusCode, errorResponse.Message, errorResponse.ErrorCode)

	} else {

		deliveryResponse := types.EndServiceDeliveryResponse{}

		err = json.Unmarshal(byteResp, &deliveryResponse)

		if err != nil {

			return types.EndServiceDeliveryResponse{}, err
		}

		return deliveryResponse, nil
	}
}
