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
	connection.Params = map[string]interface{}{}
	connection.ProcessRegex = "."
	connection.PortRange = "1234"
	var actions []Action

	var action Action
	action.TriggerRegex = ".+asdf.+"
	action.Target = "remote"
	action.TargetParams = map[string]interface {}{"address":"0.0.0.0:1234"}
	actions = append(actions, action)

	action.TriggerRegex = ".+abc.+"
	action.Target = "reply"
	action.TargetParams = map[string]interface {}{"message":"{0-10}asdf"}
	actions = append(actions, action)

	connection.Actions = actions
	connections = append(connections, connection)

	connection.ConnectionType = "ssl"
	connection.Params = map[string]interface {}{"key-path":"/tmp/key", "cert-path":"/tmp/cert"}
	connection.ProcessRegex = "."
	connection.PortRange = "100-200"

	var secondActions []Action
	action.TriggerRegex = "."
	action.Target = "remote"
	action.TargetParams = map[string]interface {}{"address":"127.0.0.1:1758"}
	secondActions = append(secondActions, action)

	connection.Actions = secondActions
	connections = append(connections, connection)

	configTesting.Connections = connections

	SetConfigPath("examples/config.json")
	if !assert.Equal(t, &configTesting, GetConfig()) {
		t.Errorf("Config didn't match expected config")
	}
}
