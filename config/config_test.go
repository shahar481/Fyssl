package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	var configTesting Config
	configTesting.LogPath = "/tmp"
	var connections []Connection

	var connection Connection
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
	action.TriggerRegex = ".+asdf.+"
	action.Target = "remote"
	action.TargetParams = map[string]interface {}{
		"listener-socket-address":"0.0.0.0:12345",
		"sender-socket-address": "0.0.0.0:12346",
	}
	actions = append(actions, action)

	action.TriggerRegex = ".+abc.+"
	action.Target = "reply"
	action.TargetParams = map[string]interface {}{"message":"{0-10}asdf"}
	actions = append(actions, action)

	connection.Actions = actions
	connections = append(connections, connection)

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
	action.TriggerRegex = "."
	action.Target = "remote"
	action.TargetParams = map[string]interface {}{
		"listener-socket-address":"127.0.0.1:1758",
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
