package configuration

import (
	"fmt"
	"strconv"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/wpwerrors"
)

// Item a configuration item has a Key to identify the item by and a value. Value is stored as a string
// but there are convenience functions to read it as other types
type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ReadInt a convenience function to attempt to return the items value as an int
func (item Item) ReadInt() (int, error) {

	i, err := strconv.Atoi(item.Value)

	if err != nil {

		return 0, wpwerrors.GetError(wpwerrors.CONVERT_VALUE, fmt.Sprintf("item.Value = \"%s\"", item.Value), err)
	}
	return i, nil
}

// ReadBool a convenience function to attempt to return the items value as a bool type
func (item Item) ReadBool() (bool, error) {

	return strconv.ParseBool(item.Value)
}
