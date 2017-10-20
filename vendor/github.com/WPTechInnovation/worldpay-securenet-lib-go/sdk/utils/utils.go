package utils

import (
	"fmt"
	"net/url"
	"reflect"
)

// StructToMap convert a struct to a map. Mapping the field name to the field value
// Only tested with string field types
// v interface{} (required) is the struct to be converted - only field values of type string have been tested
// nameTag (optional) the tag name, whos value is to be used as an alternate name instead of field name.
// errOnEmpty should the function cause an error as soon as it encounters an empty field value. This precedence over
// skipEmtpy
// skipEmtpy should the function add empty fields to the map or skip over them
func StructToMap(inputStruct interface{}, nameTag string, errOnEmpty bool, skipEmtpy bool) (map[string]string, error) {

	result := make(map[string]string, 0)

	rf := reflect.ValueOf(inputStruct)

	for i := 0; i < rf.NumField(); i++ {

		var fName string

		if nameTag == "" {
			fName = rf.Type().Field(i).Name
		} else {

			fName = rf.Type().Field(i).Tag.Get(nameTag)
		}

		fVal := rf.Field(i)

		if fVal.Len() == 0 && errOnEmpty {

			return nil, fmt.Errorf("%s is empty", fName)
		} else if fVal.Len() == 0 && !skipEmtpy {

			result[fName] = fVal.String()
		} else if fVal.Len() > 0 {

			result[fName] = fVal.String()
		}
	}

	return result, nil
}

// EncodeURLQuery add a map of string[string] to a base url using the correct encodings
// baseURL string to appear before the '?' in the result url
// inputParams key/value pairs to be encoded and appended the base url after the '?'
func EncodeURLQuery(baseURL string, inputParams map[string]string) (string, error) {

	_url, err := url.Parse(baseURL)

	if err != nil {

		return "", err
	}

	urlParams := url.Values{}

	for k, v := range inputParams {

		urlParams.Add(k, v)
	}

	_url.RawQuery = urlParams.Encode()

	return _url.String(), nil
}

// BuildUserAgentString Genereate a meaningful user agent string
func BuildUserAgentString(osName, osVersion, osArch, langVersion, libVersion, apiVersion, lang, owner string) string {

	format := "os.name=%s;os.version=%s;os.arch=%s;lang.version=%s;lib.version=%s;api.version=%s;lang=%s;owner=%s"

	return fmt.Sprintf(format, osName, osVersion, osArch, langVersion, libVersion, apiVersion, lang, owner)
}
