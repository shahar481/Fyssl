package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"slogger"
	"sync"
)

var configPath = ".config.json"
var config *Config
var once sync.Once

// GetConfig returns the config struct, will only read the config file once.
func GetConfig() *Config {
	once.Do(func() {
		config = readConfigFile()
		initConnections()
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

func initConnections() {
	for index, _ := range config.Connections {
		config.Connections[index].Id = index
		if len(config.Connections[index].ListenAddress) == 0 {
			config.Connections[index].RandomizeListenAddress = true
		} else {
			config.Connections[index].RandomizeListenAddress = false
		}
		initActions(&config.Connections[index])
	}
}

func initActions(connection *Connection) {
	for index, _ := range connection.Actions {
		compiled, err := regexp.Compile(connection.Actions[index].TriggerRegex)
		if err != nil {
			slogger.Error(fmt.Sprintf("Error compiling trigger number %d in connection-%s",index, connection.Name))
			panic("Error compiling regex")
		}
		connection.Actions[index].CompiledTrigger = compiled
		connection.Actions[index].Id = index
	}
}

func initRemote(action *Action) {
	if action.Target == "remote" {
		action.TargetParams.(map[string]interface{})["sockets"] = make(map[string]net.Conn)
	}
}

// SetConfigPath sets the config path to read the config from.
// Don't use it after you have used GetConfig once.
func SetConfigPath(c string) {
	configPath = c
}