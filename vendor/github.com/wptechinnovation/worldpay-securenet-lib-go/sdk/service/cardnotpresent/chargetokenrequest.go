package cardnotpresent

import "github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/types"

// ChargeTokenRequest represents a request to charge a token where the card is not present
type ChargeTokenRequest struct {
	Amount               float32                     `json:"amount"`
	PaymentVaultToken    *types.PaymentVaultToken    `json:"paymentVaultToken"`
	DeveloperApplication *types.DeveloperApplication `json:"developerApplication"`
	ExtendedInformation  *types.ExtendedInformation  `json:"extendedInformation"`
}
