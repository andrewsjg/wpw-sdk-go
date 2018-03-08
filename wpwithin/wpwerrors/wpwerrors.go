package wpwerrors

import (
	"fmt"
)

const (
	GENERAL = "GENERAL"
	NET     = "NET"
	MEMORY  = "MEMORY"
	DATA    = "DATA"
	MATH    = "MATH"
	JSON    = "JSON"
	HTTP    = "HTTP"
	CONFIG  = "CONFIG"
	IO      = "IO"
	CORE    = "CORE"
	HTE     = "HTE"
	PSP     = "PSP"
	UTILS   = "UTILS"
	WPW     = "WPW"
)

type ErrorID int

const (
	UNKNOWN ErrorID = iota
	NO_DATA
	DIVISION_BY_ZERO
	IO_OPEN
	IO_READ
	IO_WRITE
	IO_SYNC
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
	HTTP_REQUEST_POST
	HTTP_REQUEST_DO
	EMPTYMERCHANTCLIENTKEY
	EMPTYMERCHANTSERVICEKEY
	ORDER_EXISTS
	ORDER_NOTFOUND
	NOT_IMPLEMENTED
	EMPTY_CLIENTTOKEN
	EMPTY_ORDERDESCRIPTION
	EMPTY_CUSTOMERORDERCODE
	PAYMENT_FAILED
	POST_FAILED
	SNCLIENT_CREATEFAILED
	LISTEN_FAILED
	NEWUDPCOMMERR
	SCAN4SERVICESERR
	LISTOFINTFS
	LISTOFADDRS
	DEVICENOTCONN
	CALCNETMASK
	UUIDGEN
	EMPTYNAME
	EMPTYDESCRIPTION
	SDKFACTORYCREATE
	CORECREATION
	FACTORYGETDEVICE
	FACTORYGETORDERMANAGER
	FACTORYGETSVCBROADCAST
	FACTORYGETSVCSCANNER
	SERVICEEXISTS
	FACTORYGETPSPCLIENT
	FACTORYGETHTECLIENTHTTP
	HTENEWCLIENT
	PSPMERCHANTCREATE
	HTENEWCREDENTIAL
	FACTORYGETHTE
	SCANFORSERVICES
	SCANFORSERVICE
	HTECLIENTGETPRICES
	PRICESNIL
	PSPGETTOKEN
	DISCOVERSERVICES
	NILSERVISERESPONSE
	STARTDELIVERYFAILED
	ENDDELIVERYFAILED
)

type WpwError struct {
	ID      string
	Type    string // ErrorType
	Message string
}

var errors = [...]WpwError{
	UNKNOWN:                 {"UNKNOWN", GENERAL, "unknow error"},
	NO_DATA:                 {"NO_DATA", DATA, "lack of data"},
	DIVISION_BY_ZERO:        {"DIVISION_BY_ZERO", MATH, "division by zero"},
	IO_OPEN:                 {"IO_OPEN", IO, "io open error"},
	IO_READ:                 {"IO_READ", IO, "io read error"},
	IO_WRITE:                {"IO_WRITE", IO, "io write error"},
	IO_SYNC:                 {"IO_SYNC", IO, "io sync error"},
	DECODE_JSON:             {"DECODE_JSON", JSON, "failed to decode json"},
	ENCODE_JSON:             {"ENCODE_JSON", JSON, "failed to encode json"},
	HTTP_GET:                {"HTTP_GET", HTTP, "http get error"},
	HTTP_POSTJSON:           {"HTTP_POSTJSON", HTTP, "http failed to post JSON"},
	HTTP_REQUEST_POST:       {"HTTP_REQUEST_POST", HTTP, "http request post error"},
	HTTP_REQUEST_DO:         {"HTTP_REQUEST_DO", HTTP, "http request do error"},
	WRONG_CONFIG_PATH:       {"WRONG_CONFIG_PATH", CONFIG, "wrong path to configuration data"},
	OPEN_FILE:               {"OPEN_FILE", CONFIG, "failed to open file"},
	CONVERT_VALUE:           {"CONVERT_VALUE", CONFIG, "failed to convert value"},
	UUID_FILE_READ:          {"UUID_FILE_READ", CORE, "failed to read uuid file"},
	UUID_FILE_CREATE:        {"UUID_FILE_CREATE", CORE, "failed to create uuid file"},
	UUID_FILE_SAVE:          {"UUID_FILE_SAVE", CORE, "failed to save uuid file"},
	PSP_CONFIG_NOT_SET:      {"PSP_CONFIG_NOT_SET", CORE, "PSP Config is not set"},
	PARSE_DEVELOPER_ID:      {"PARSE_DEVELOPER_ID", CORE, "failed to parse developerID"},
	PSP_UNKNOWN:             {"PSP_UNKNOWN", CORE, "unknown psp"},
	PSP_COLLECTION:          {"PSP_COLLECTION", CORE, "PSP Config collection is be set"},
	EMPTY_HOST:              {"EMPTY_HOST", HTE, "empty host"},
	PORT_RANGE:              {"PORT_RANGE", HTE, "port number out of range"},
	EMPTY_CLIENTID:          {"EMPTY_CLIENTID", HTE, "empty clientId"},
	EMPTYMERCHANTCLIENTKEY:  {"EMPTYMERCHANTCLIENTKEY", HTE, "empty MerchantClientKey"},
	EMPTYMERCHANTSERVICEKEY: {"EMPTYMERCHANTSERVICEKEY", HTE, "empty MerchantServiceKey"},
	ORDER_EXISTS:            {"ORDER_EXISTS", HTE, "order already exists"},
	ORDER_NOTFOUND:          {"ORDER_NOTFOUND", HTE, "order not found"},
	NOT_IMPLEMENTED:         {"NOT_IMPLEMENTED", PSP, "functionality not implemented"},
	EMPTY_CLIENTTOKEN:       {"EMPTY_CLIENTTOKEN", PSP, "empty clientToken"},
	EMPTY_ORDERDESCRIPTION:  {"EMPTY_ORDERDESCRIPTION", PSP, "empty orderDescription"},
	EMPTY_CUSTOMERORDERCODE: {"EMPTY_CUSTOMERORDERCODE", PSP, "empty customerOrderCode"},
	PAYMENT_FAILED:          {"PAYMENT_FAILED", PSP, "payment failed"},
	POST_FAILED:             {"POST_FAILED", PSP, "post failed"},
	SNCLIENT_CREATEFAILED:   {"SNCLIENT_CREATEFAILED", PSP, "failed to create new SN client"},
	LISTEN_FAILED:           {"LISTEN_FAILED", NET, "failed to listen"},
	NEWUDPCOMMERR:           {"NEWUDPCOMMERR", NET, "failed to create udp communicator"},
	SCAN4SERVICESERR:        {"SCAN4SERVICESERR", NET, "ScanForServices failed"},
	LISTOFINTFS:             {"LISTOFINTFS", NET, "failed to get list of interfaces"},
	LISTOFADDRS:             {"LISTOFADDRS", NET, "failed to get list of addreses"},
	DEVICENOTCONN:           {"DEVICENOTCONN", NET, "device does not appear to be network connected"},
	CALCNETMASK:             {"CALCNETMASK", NET, "unable to calculate netmask"},
	UUIDGEN:                 {"UUIDGEN", UTILS, "failed to generate random UUID"},
	EMPTYNAME:               {"EMPTYNAME", WPW, "name is empty"},
	EMPTYDESCRIPTION:        {"EMPTYDESCRIPTION", WPW, "description is empty"},
	SDKFACTORYCREATE:        {"SDKFACTORYCREATE", WPW, "unable to create SDK factory"},
	CORECREATION:            {"CORECREATION", WPW, "error creating new core.Core"},
	FACTORYGETDEVICE:        {"FACTORYGETDEVICE", WPW, "unable to get device"},
	FACTORYGETORDERMANAGER:  {"FACTORYGETORDERMANAGER", WPW, "unable to get order manager"},
	FACTORYGETSVCBROADCAST:  {"FACTORYGETSVCBROADCAST", WPW, "unable to get service broadcaster"},
	FACTORYGETSVCSCANNER:    {"FACTORYGETSVCSCANNER", WPW, "unable to get service scanner"},
	SERVICEEXISTS:           {"SERVICEEXISTS", WPW, "service with that id already exists"},
	FACTORYGETPSPCLIENT:     {"FACTORYGETPSPCLIENT", WPW, "unable to get PSP client"},
	FACTORYGETHTECLIENTHTTP: {"FACTORYGETHTECLIENTHTTP", WPW, "unable to get HTE client HTTP"},
	HTENEWCLIENT:            {"HTENEWCLIENT", WPW, "unable to create HTE client"},
	PSPMERCHANTCREATE:       {"PSPMERCHANTCREATE", WPW, "unable to create psp"},
	HTENEWCREDENTIAL:        {"HTENEWCREDENTIAL", WPW, "unable to create new HTE credential"},
	FACTORYGETHTE:           {"FACTORYGETHTE", WPW, "unable to get HTE"},
	SCANFORSERVICES:         {"SCANFORSERVICES", WPW, "scan for services failed"},
	SCANFORSERVICE:          {"SCANFORSERVICE", WPW, "scan for service failed"},
	HTECLIENTGETPRICES:      {"HTECLIENTGETPRICES", WPW, "unable to get HTE client prices"},
	PRICESNIL:               {"PRICESNIL", WPW, "response with nil prices"},
	PSPGETTOKEN:             {"PSPGETTOKEN", WPW, "failed to get token"},
	DISCOVERSERVICES:        {"DISCOVERSERVICES", WPW, "failed to discover services"},
	NILSERVISERESPONSE:      {"NILSERVISERESPONSE", WPW, "discover services is successful but serviceResponse is nil"},
	STARTDELIVERYFAILED:     {"STARTDELIVERYFAILED", WPW, "start delivery failed"},
	ENDDELIVERYFAILED:       {"ENDDELIVERYFAILED", WPW, "end delivery failed"},
}

// Error return formated error for specified error (e).
func (e WpwError) Error() error {
	return fmt.Errorf("%v, %v, %v", e.ID, e.Type, e.Message)
}

const sep = ", "
const openBracket = " :{"
const closeBracket = "}"
const unsupportedTypeError = "<unsupported type of passed data>"

// GetError returns error for specified eid with additional data (if any).
func GetError(eid ErrorID, additionalData ...interface{}) error {
	if len(additionalData) == 0 {
		return errors[eid].Error()
	}

	var ret string
	for _, v := range additionalData {
		switch v := v.(type) {
		case string:
			if ret == "" {
				ret = sep + v
			} else {
				ret = ret + sep + v
			}
		case error:
			// if the additionalData contains other value of type error
			// then it will be closed in brackets
			if ret == "" {
				ret = openBracket + v.Error() + closeBracket
			} else {
				ret = ret + openBracket + v.Error() + closeBracket
			}
		default:
			if ret == "" {
				ret = unsupportedTypeError
			} else {
				ret = ret + sep + unsupportedTypeError
			}
		}
	}
	return fmt.Errorf("%v%v", errors[eid].Error(), ret)
}
