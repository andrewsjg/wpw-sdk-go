package errors

import (
	"fmt"
	"strings"
	"time"
)

type ErrorId int

const (
	UNKNOWN = 0 + iota
	NO_DATA
)

type ErrorType int

const (
	GENERAL = "GENERAL"
	NET     = "NET"
	MEMORY  = "MEMEORY"
	DATA    = "DATA"
)

type WpwError struct {
	ID      string
	Type    string // ErrorType
	Message string
}

var errors = [...]WpwError{
	UNKNOWN: {"UNKNOWN", GENERAL, "unknow error"},
	// add new error below
	NO_DATA: {"NO_DATA", DATA, "lack of data"},
}

// Error return formated error for specified error (e).
func (e WpwError) Error() error {
	t := time.Now()
	formatedTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d-00:00", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return fmt.Errorf("%v: %v, %v, %s", formatedTime, e.ID, e.Type, e.Message)
}

// GetError returns error for specified eid with additional data (if any) separated with comma.
func GetError(eid ErrorId, additionalData ...string) error {
	if len(additionalData) == 0 {
		return errors[eid].Error()
	}
	return fmt.Errorf("%v, %v", errors[eid].Error(), strings.Join(additionalData, ", "))
}
