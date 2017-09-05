package onlineworldpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay/types"
	wpwithin_types "github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
)

// OnlineWorldpay implementation of PSP
type OnlineWorldpay struct {
	MerchantClientKey  string
	MerchantServiceKey string
	apiEndpoint        string
}

// NewMerchant create a new clint with a merchant context
func NewMerchant(merchantClientKey, merchantServiceKey, apiEndpoint string) (psp.PSP, error) {

	result := &OnlineWorldpay{
		MerchantClientKey:  merchantClientKey,
		MerchantServiceKey: merchantServiceKey,
		apiEndpoint:        apiEndpoint,
	}

	return result, nil
}

// NewClient create a new client - no merchant information is required
func NewClient(apiEndpoint string) (psp.PSP, error) {

	result := &OnlineWorldpay{
		apiEndpoint: apiEndpoint,
	}

	return result, nil
}

// GetToken by passing a card credentials and a clientKey, Worldpay returns a token representing
// those payment credentials
func (owp *OnlineWorldpay) GetToken(hceCredentials *wpwithin_types.HCECard, clientKey string, reusableToken bool) (string, error) {

	log.Debug("begin onlineworldpay.GetToken()")

	if reusableToken {
		// TODO: CH - support reusable token by storing the value (along with merchant client key so link to a merchant) within the car so that token can be re-used if present, or created if not
		return "", errors.New("Reusable token support not implemented")
	}

	paymentMethod := types.TokenRequestPaymentMethod{
		Name:        fmt.Sprintf("%s %s", hceCredentials.FirstName, hceCredentials.LastName),
		ExpiryMonth: hceCredentials.ExpMonth,
		ExpiryYear:  hceCredentials.ExpYear,
		CardNumber:  hceCredentials.CardNumber,
		Type:        hceCredentials.Type,
		Cvc:         hceCredentials.Cvc,
		StartMonth:  nil,
		StartYear:   nil,
	}

	tokenRequest := types.TokenRequest{
		Reusable:      reusableToken,
		PaymentMethod: paymentMethod,
		ClientKey:     clientKey,
	}

	log.Debug("Attempt to marshal token request to JSON")
	bJSON, err := json.Marshal(tokenRequest)

	if err != nil {

		return "", err
	}

	log.WithField("Did marshal TokenRequest JSON", string(bJSON)).Debug("POST Request Token.")

	reqURL := fmt.Sprintf("%s/tokens", owp.apiEndpoint)

	var tokenResponse types.TokenResponse

	log.WithFields(log.Fields{"Url": reqURL,
		"RequestJSON": string(bJSON)}).Debug("Sending Token POST request.")
	err = post(reqURL, bJSON, make(map[string]string, 0), &tokenResponse)

	if err != nil {

		log.WithField("Error", err).Error("Error POSTing")
	}

	return tokenResponse.Token, err
}

// MakePayment make a payment
func (owp *OnlineWorldpay) MakePayment(amount int, currencyCode, clientToken, orderDescription, customerOrderCode string) (string, error) {

	log.WithFields(log.Fields{"Amount": strconv.Itoa(amount), "CurrencyCode": currencyCode, "ClientToken": clientToken,
		"OrderDescription": orderDescription, "CustomerOrderCode": customerOrderCode}).Debug("Begin OWP MakePayment")

	if clientToken == "" {

		return "", errors.New("clientToken cannot be empty")
	}
	if orderDescription == "" {

		return "", errors.New("orderDescription cannot be empty")
	}
	if customerOrderCode == "" {

		return "", errors.New("customerOrderCode cannot be empty")
	}

	orderRequest := types.OrderRequest{

		Token:             clientToken,
		Amount:            amount,
		CurrencyCode:      currencyCode,
		OrderDescription:  orderDescription,
		CustomerOrderCode: customerOrderCode,
	}

	bJSON, err := json.Marshal(orderRequest)

	if err != nil {

		return "", err
	}

	log.WithField("JSON", string(bJSON)).Debug("JSON form of OrderRequest object.")

	reqURL := fmt.Sprintf("%s/orders", owp.apiEndpoint)

	log.WithFields(log.Fields{"Request URL": reqURL, "MerchantSvcKey": owp.MerchantServiceKey}).Debug("Using OWP parameters")

	var orderResponse types.OrderResponse

	headers := make(map[string]string, 0)

	headers["Authorization"] = owp.MerchantServiceKey

	log.WithFields(log.Fields{"Url": reqURL,
		"RequestJSON": string(bJSON)}).Debug("Sending Order POST request.")

	err = post(reqURL, bJSON, headers, &orderResponse)

	if err != nil {

		return "", err
	}

	if strings.EqualFold(orderResponse.PaymentStatus, "SUCCESS") {

		return orderResponse.OrderCode, nil
	}

	return "", fmt.Errorf("Payment failed for customer order %s ", orderResponse.CustomerOrderCode)
}

func post(url string, requestBody []byte, headers map[string]string, v interface{}) error {

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	for k, v := range headers {

		request.Header.Set(k, v)
	}

	if err != nil {

		return err
	}

	// TODO: CH Add a http client as a dependency during construction to aid testing
	client := &http.Client{}

	resp, err := client.Do(request)

	if err != nil {

		return err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		return err
	}

	log.WithField("Code", resp.StatusCode).Debug("Response status code")
	log.Debug(fmt.Sprintf("Response body: %s", string(respBody)))

	if resp.StatusCode == HTTPOK {

		log.Debug("POST response code = HTTPOK, attempting to unmarshal and return")
		return json.Unmarshal(respBody, &v)
	}

	log.WithField("HTTP Code", resp.StatusCode).Debug("POST response code != HTTPOK. Attepting to parse error response message")

	wpErr := types.ErrorResponse{}

	err = json.Unmarshal(respBody, &wpErr)
	if err != nil {

		return err
	}

	log.WithFields(log.Fields{"Message": wpErr.Message, "Description": wpErr.Description, "CustomCode": wpErr.CustomCode, "HTTP Status Code": wpErr.HTTPStatusCode, "HelpUrl": wpErr.ErrorHelpURL}).Debug("** POST Response")

	return fmt.Errorf("HTTP Status: %d - CustomCode: %s - Message: %s", wpErr.HTTPStatusCode, wpErr.CustomCode, wpErr.Message)
}
