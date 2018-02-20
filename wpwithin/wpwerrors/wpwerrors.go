package wpwerrors

import (
	"fmt"
	"strings"
)

const (
	GENERAL = "GENERAL"
	NET     = "NET"
	MEMORY  = "MEMORY"
	DATA    = "DATA"
	MATH    = "MATH"
	JSON    = "JSON"
	CONFIG  = "CONFIG"
	IO      = "IO"
	CORE    = "CORE"
	HTE     = "HTE"
)

type ErrorID int

const (
	UNKNOWN ErrorID = iota
	NO_DATA
	DIVISION_BY_ZERO
	WRONG_CONFIG_PATH
	OPEN_FILE
	DECODE_JSON
	ENCODE_JSON
	CONVERT_VALUE
	UUID_FILE_READ
	UUID_FILE_CREATE
	UUID_FILE_SAVE
	PSP_CONFIG_NOT_SET
	PARSE_DEVELOPER_ID
	PSP_UNKNOWN
	PSP_COLLECTION
	EMPTY_HOST
	PORT_RANGE
	EMPTY_CLIENTID
	HTTP_GET
	HTTP_POSTJSON
)

type WpwError struct {
	ID      string
	Type    string // ErrorType
	Message string
}

var errors = [...]WpwError{
	UNKNOWN:            {"UNKNOWN", GENERAL, "unknow error"},
	NO_DATA:            {"NO_DATA", DATA, "lack of data"},
	DIVISION_BY_ZERO:   {"DIVISION_BY_ZERO", MATH, "division by zero"},
	DECODE_JSON:        {"DECODE_JSON", JSON, "failed to decode json"},
	ENCODE_JSON:        {"ENCODE_JSON", JSON, "failed to encode json"},
	WRONG_CONFIG_PATH:  {"WRONG_CONFIG_PATH", CONFIG, "wrong path to configuration data"},
	OPEN_FILE:          {"OPEN_FILE", CONFIG, "failed to open file"},
	CONVERT_VALUE:      {"CONVERT_VALUE", CONFIG, "failed to convert value"},
	UUID_FILE_READ:     {"UUID_FILE_READ", CORE, "failed to read uuid file"},
	UUID_FILE_CREATE:   {"UUID_FILE_CREATE", CORE, "failed to create uuid file"},
	UUID_FILE_SAVE:     {"UUID_FILE_SAVE", CORE, "failed to save uuid file"},
	PSP_CONFIG_NOT_SET: {"PSP_CONFIG_NOT_SET", CORE, "PSP Config is not set"},
	PARSE_DEVELOPER_ID: {"PARSE_DEVELOPER_ID", CORE, "failed to parse developerID"},
	PSP_UNKNOWN:        {"PSP_UNKNOWN", CORE, "unknown psp"},
	PSP_COLLECTION:     {"PSP_COLLECTION", CORE, "PSP Config collection is be set"},
	EMPTY_HOST:         {"EMPTY_HOST", HTE, "empty host"},
	PORT_RANGE:         {"PORT_RANGE", HTE, "port number out of range"},
	EMPTY_CLIENTID:     {"EMPTY_CLIENTID", HTE, "empty clientId"},
	HTTP_GET:           {"HTTP_GET", HTE, "http client get error"},
	HTTP_POSTJSON:      {"HTTP_POSTJSON", HTE, "http client failed to post JSON"},
}

// Error return formated error for specified error (e).
func (e WpwError) Error() error {
	return fmt.Errorf("%v, %v, %v", e.ID, e.Type, e.Message)
}

// GetError returns error for specified eid with additional data (if any) separated with comma.
func GetError(eid ErrorID, additionalData ...string) error {
	if len(additionalData) == 0 {
		return errors[eid].Error()
	}
	return fmt.Errorf("%v, %v", errors[eid].Error(), strings.Join(additionalData, ", "))
}
