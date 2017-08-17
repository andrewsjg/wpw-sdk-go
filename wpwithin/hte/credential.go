package hte

import "errors"

// Credential Merchant HTE Credentials
type Credential struct {
	MerchantClientKey  string
	MerchantServiceKey string
}

// NewHTECredential create a new HTE Credential. All parameters are required
func NewHTECredential(MerchantClientKey, MerchantServiceKey string) (*Credential, error) {

	if MerchantClientKey == "" {

		return nil, errors.New("MerchantClientKey required, cannot be empty")

	} else if MerchantServiceKey == "" {

		return nil, errors.New("MerchantServiceKey required, cannot be empty")
	}

	result := &Credential{
		MerchantClientKey:  MerchantClientKey,
		MerchantServiceKey: MerchantServiceKey,
	}

	return result, nil
}
