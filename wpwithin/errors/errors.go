package errors

import (
	"fmt"
	"strings"
)

type ErrorId int

const (
	UNKNOWN ErrorId = iota
	NO_DATA
	DIVISION_BY_ZERO
)

const (
	GENERAL = "GENERAL"
	NET     = "NET"
	MEMORY  = "MEMEORY"
	DATA    = "DATA"
	MATH    = "MATH"
)

type WpwError struct {
	ID      string
	Type    string // ErrorType
	Message string
}

var errors = [...]WpwError{
	UNKNOWN:          {"UNKNOWN", GENERAL, "unknow error"},
	NO_DATA:          {"NO_DATA", DATA, "lack of data"},
	DIVISION_BY_ZERO: {"DIVISION_BY_ZERO", MATH, "division by zero"},
}

// Error return formated error for specified error (e).
func (e WpwError) Error() error {
	return fmt.Errorf("%v, %v, %v", e.ID, e.Type, e.Message)
}

// GetError returns error for specified eid with additional data (if any) separated with comma.
func GetError(eid ErrorId, additionalData ...string) error {
	if len(additionalData) == 0 {
		return errors[eid].Error()
	}
	return fmt.Errorf("%v, %v", errors[eid].Error(), strings.Join(additionalData, ", "))
}
