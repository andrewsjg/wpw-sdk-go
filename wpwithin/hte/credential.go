package hte

import (
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/wpwerrors"
)

// Credential Merchant HTE Credentials
type Credential struct {
	MerchantClientKey  string
	MerchantServiceKey string
}

// NewHTECredential create a new HTE Credential. All parameters are required
func NewHTECredential(MerchantClientKey, MerchantServiceKey string) (*Credential, error) {

	if MerchantClientKey == "" {
		return nil, wpwerrors.GetError(wpwerrors.EMPTYMERCHANTCLIENTKEY)
	}
	if MerchantServiceKey == "" {
		return nil, wpwerrors.GetError(wpwerrors.EMPTYMERCHANTSERVICEKEY)
	}

	result := &Credential{
		MerchantClientKey:  MerchantClientKey,
		MerchantServiceKey: MerchantServiceKey,
	}

	return result, nil
}
