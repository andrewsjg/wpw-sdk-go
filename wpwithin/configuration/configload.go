package configuration

import (
	"encoding/json"
	"errors"
	"os"
)

// Configuration a configuration helper. Will parse a configuration JSON file into memory and store in a map[string]Item
// For each access to configuration values
type Configuration struct {
	items map[string]Item
}

// Load using a specified path, parse the file as a map of Item
func Load(configPath string) (Configuration, error) {

	if configPath == "" {

		return Configuration{}, errors.New("configPath is not set")
	}

	file, err := os.Open(configPath)

	if err != nil {

		return Configuration{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var tmpConfig []Item
	err = decoder.Decode(&tmpConfig)

	if err != nil {

		return Configuration{}, err
	}
	result := Configuration{}
	result.items = make(map[string]Item, 0)

	if len(tmpConfig) > 0 {

		for _, v := range tmpConfig {

			result.items[v.Key] = v
		}
	}

	return result, nil
}

// GetValue a convenience function to read the configuration item for a specific key
func (config Configuration) GetValue(key string) Item {

	return config.items[key]
}

// GetItems returns a map representative of the loaded config file. The index value in the map is the key value
// in each config item
func (config Configuration) GetItems() map[string]Item {

	return config.items
}
