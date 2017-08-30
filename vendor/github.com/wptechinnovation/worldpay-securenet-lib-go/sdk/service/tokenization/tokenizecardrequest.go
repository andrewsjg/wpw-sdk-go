package tokenization

import "github.com/wptechinnovation/worldpay-securenet-lib-go/sdk/types"

// TokenizeCardRequest represents a request to convert a payment card to a SecureNet token
type TokenizeCardRequest struct {
	Card                 *types.Card                 `json:"card"`
	PublicKey            string                      `json:"publicKey"`
	DeveloperApplication *types.DeveloperApplication `json:"developerApplication"`
	AddToVault           bool                        `json:"addToVault"`
	CustomerID           *string                     `json:"customerId,omitempty"`
}
