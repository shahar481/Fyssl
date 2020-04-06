package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var configPath = ".config.json"
var config *Config
var once sync.Once

// GetConfig returns the config struct, will only read the config file once.
func GetConfig() *Config {
	once.Do(func() {
		config = readConfigFile()
	})
	return config
}

func readConfigFile() *Config {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config *Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}

// SetConfigPath sets the config path to read the config from.
// Don't use it after you have used GetConfig once.
func SetConfigPath(c string) {
	configPath = c
}