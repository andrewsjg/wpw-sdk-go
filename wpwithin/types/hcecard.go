package types

import (
	"fmt"
	"strings"
)

// HCECard represents details of a payment card
type HCECard struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	ExpMonth   int32  `json:"expMonth"`
	ExpYear    int32  `json:"expYear"`
	CardNumber string `json:"cardNumber"`
	Type       string `json:"type"`
	Cvc        string `json:"cvc"`
}

// ObfuscateCardNumber obfuscates the card number with the exception of last four digits
// i.e. 1234567890987654 becomes ************7654
func (card HCECard) ObfuscateCardNumber() string {

	pan := card.CardNumber

	if len(pan) <= 4 {

		return pan
	}

	var result = make([]string, 0)

	for i := 1; i < len(pan)-4; i++ {

		result = append(result, "*")
	}
	result = append(result, pan[len(pan)-4:len(pan)])

	return strings.Join(result, "")
}

func (card HCECard) String() string {

	return fmt.Sprintf("fname=%s;lname=%s;expm=%d;expy%d;num=%s;cvc=%s;type=%s;", card.FirstName, card.LastName, card.ExpMonth, card.ExpYear, card.ObfuscateCardNumber(), card.Cvc, card.Type)
}
