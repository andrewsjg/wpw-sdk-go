package types

type PspConfig struct {
	PspName            string `json:"psp_name"`
	ApiEndpoint        string `json:"api_endpoint"`
	HtePublicKey       string `json:"hte_public_key,omitempty"`
	HtePrivateKey      string `json:"hte_private_key,omitempty"`
	MerchantClientKey  string `json:"merchant_client_key,omitempty"`
	MerchantServiceKey string `json:"merchant_service_key,omitempty"`
	DeveloperId        string `json:"developer_id,omitempty"`
	AppVersion         string `json:"app_version,omitempty"`
	PublicKey          string `json:"public_key,omitempty"`
	SecureKey          string `json:"secure_key,omitempty"`
	SecureNetId        string `json:"secure_net_id,omitempty"`
}
