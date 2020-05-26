package config

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestGetConfig(t *testing.T) {
	var configTesting Config
	configTesting.LogPath = "/tmp"
	var connections []Connection

	var connection Connection
	connection.Name = "tcp-test"
	connection.ConnectionType = "tcp"
	connection.ListenAddress = "0.0.0.0:4576"
	connection.ConnectAddress = "127.0.0.1:1234"
	connection.Params = map[string]interface{}{}
	connection.Modules = map[string]interface{}{
		"catch": map[string]interface{} {
			"ports": "1234",
			"process":".",
		},
	}
	var actions []Action

	var action Action
	action.Name = "python-script-processor"
	action.TriggerRegex = ".+asdf.+"
	action.DumpLogPackets = false
	compiled, err := regexp.Compile(action.TriggerRegex)
	if err != nil {
		panic("Error compiling regex")
	}
	action.CompiledTrigger = compiled
	action.Target = "remote"
	action.TargetParams = map[string]interface {}{
		"listener-socket-address":"0.0.0.0:12345",
		"sender-socket-address": "0.0.0.0:12346",
	}
	actions = append(actions, action)

	action.Name = "asdf-replier"
	action.TriggerRegex = ".+abc.+"
	compiled, err = regexp.Compile(action.TriggerRegex)
	if err != nil {
		panic("Error compiling regex")
	}
	action.CompiledTrigger = compiled
	action.DumpLogPackets = true
	action.Target = "reply"
	action.TargetParams = map[string]interface {}{"message":"{0-10}asdf"}
	actions = append(actions, action)

	connection.Actions = actions
	connections = append(connections, connection)

	connection.Name = "ssl-test"
	connection.ConnectionType = "ssl"
	connection.ListenAddress = "0.0.0.0:15689"
	connection.ConnectAddress = "127.0.0.1:1876"
	connection.Params = map[string]interface {}{
		"key-path":"/tmp/key",
		"cert-path":"/tmp/cert",
	}
	connection.Modules = map[string]interface{} {
		"catch": map[string]interface {} {
			"process": ".",
			"ports": "100-200",
		},
	}

	var secondActions []Action
	action.Name = "ruby-script-processor"
	action.TriggerRegex = "."
	action.DumpLogPackets = false
	compiled, err = regexp.Compile(action.TriggerRegex)
	if err != nil {
		panic("Error compiling regex")
	}
	action.CompiledTrigger = compiled
	action.Target = "remote"
	action.TargetParams = map[string]interface {}{
		"listener-socket-address":"127.0.0.1:17586",
		"sender-socket-address": "",
	}
	secondActions = append(secondActions, action)

	connection.Actions = secondActions
	connections = append(connections, connection)

	configTesting.Connections = connections

	SetConfigPath("examples/config.json")
	if !assert.Equal(t, &configTesting, GetConfig()) {
		t.Errorf("Config didn't match expected config")
	}
}
