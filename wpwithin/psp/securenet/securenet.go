package securenet

import (
	"fmt"
	"strconv"

	snclient "github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/client"
	"github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/service/cardnotpresent"
	"github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/service/tokenization"
	sntypes "github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/types"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/utils"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/wpwerrors"
)

// SecureNet supports processing of payments using the Secure Net service
type SecureNet struct {
	publicKey    string
	developerApp sntypes.DeveloperApplication
	snclient     snclient.Client
}

// NewSecureNetMerchant creates a new instance of the SecureNet client with a merchant context
func NewSecureNetMerchant(snID, secureKey, publicKey, apiEndpoint, appVersion string, developerID int32, proxy string) (psp.PSP, error) {

	result := SecureNet{}

	result.developerApp = sntypes.DeveloperApplication{}
	result.developerApp.DeveloperID = developerID
	result.developerApp.Version = appVersion
	result.publicKey = publicKey
	_snclient, err := snclient.New(apiEndpoint, appVersion, snID, secureKey, proxy)

	if err != nil {

		return nil, wpwerrors.GetError(wpwerrors.SNCLIENT_CREATEFAILED, err.Error())
	}

	result.snclient = _snclient

	return &result, nil

}

// NewSecureNetConsumer creates a new instance of the SecureNet client with a consumer context
func NewSecureNetConsumer(apiEndpoint, appVersion string, developerID int32, proxy string) (psp.PSP, error) {

	result := SecureNet{}

	result.developerApp = sntypes.DeveloperApplication{}
	result.developerApp.DeveloperID = developerID
	result.developerApp.Version = appVersion

	_snclient, err := snclient.New(apiEndpoint, appVersion, "", "", proxy)

	if err != nil {

		return nil, wpwerrors.GetError(wpwerrors.SNCLIENT_CREATEFAILED, err.Error())
	}

	result.snclient = _snclient

	return &result, nil
}

// GetToken converts a payment credential into a token
func (sn *SecureNet) GetToken(hceCredentials *types.HCECard, clientKey string, reusableToken bool) (string, error) {

	// Create a card with test data
	card := sntypes.Card{}
	card.Number = hceCredentials.CardNumber
	card.ExpirationDate = fmt.Sprintf("%d/%d", hceCredentials.ExpMonth, hceCredentials.ExpYear)

	// Attempt to tokenize the card
	reqTokenize := tokenization.TokenizeCardRequest{}
	reqTokenize.Card = &card
	reqTokenize.AddToVault = reusableToken
	reqTokenize.PublicKey = clientKey
	reqTokenize.DeveloperApplication = &sn.developerApp

	respTokenize, err := sn.snclient.TokenizationService().TokenizeCard(&reqTokenize)

	var tokenStr string

	if err == nil {

		tokenStr = respTokenize.Token
	}

	return tokenStr, err
}

// MakePayment charges a payment token as if it was a payment credential
func (sn *SecureNet) MakePayment(iAmount int, currencyCode, clientToken, orderDescription, customerOrderCode string) (string, error) {

	// Need to convert the amount which is minor units (integer) into a value that
	// represents major units (float)
	minor := float32((iAmount % 100)) / 100
	major := float32(iAmount / 100)
	converted := utils.ToFixed(float64(major+minor), 2)

	reqChargeToken := cardnotpresent.ChargeTokenRequest{}
	reqChargeToken.Amount = float32(converted)
	reqChargeToken.DeveloperApplication = &sn.developerApp
	reqChargeToken.PaymentVaultToken = &sntypes.PaymentVaultToken{}
	reqChargeToken.PaymentVaultToken.CustomerID = customerOrderCode
	reqChargeToken.PaymentVaultToken.PaymentMethodID = clientToken
	reqChargeToken.PaymentVaultToken.PaymentType = sntypes.CreditCard
	reqChargeToken.PaymentVaultToken.PublicKey = sn.publicKey

	respCharge, err := sn.snclient.CardNotPresentService().ChargeUsingToken(&reqChargeToken)

	var transactionID string

	if err == nil {

		transactionID = strconv.Itoa(respCharge.Transaction.TransactionID)
	}

	return transactionID, err
}
