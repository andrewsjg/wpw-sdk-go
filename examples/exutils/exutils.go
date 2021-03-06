package exutils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
)

type Config struct {
	LogFileName   string          `json:"logFileName,omitempty"`
	DeviceName    string          `json:"deviceName,omitempty"`
	InterfaceAddr string          `json:"interfaceAddr,omitempty"`
	HceCard       types.HCECard   `json:"hceCard,omitempty"`
	PspConfig     types.PspConfig `json:"pspConfig"`
}

func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}
