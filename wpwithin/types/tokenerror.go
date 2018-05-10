package types

import "github.com/wptechinnovation/wpw-sdk-go/wpwithin/types/errors"

type TokenError struct {
	ErrorType errors.ErrorType
	DetailMsg string
}
