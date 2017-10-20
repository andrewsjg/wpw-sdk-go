package psp

import (
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
)

// Common configuration map keys
const (
	// CfgPSPName Common configuration key, name of PSP
	CfgPSPName string = "psp_name"
	// CfgHTEPublicKey ..
	CfgHTEPublicKey string = "hte_public_key"
	// CfgHTEPrivateKey ..
	CfgHTEPrivateKey string = "hte_private_key"
)

// PSP defines functions for making payments
type PSP interface {
	GetToken(hceCredentials *types.HCECard, clientKey string, reusableToken bool) (string, error)
	MakePayment(amount int, currencyCode, clientToken, orderDescription, customerOrderCode string) (string, error)
}
