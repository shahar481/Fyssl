package config

import (
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	var configTesting Config
	configTesting.LogPath = "/tmp"
	var connections []Connection

	var connection Connection
	connection.ConnectionType = "tcp"
	connection.Params = map[string]string{}
	connection.ProcessRegex = "."
	connection.PortRange = "1234"
	var actions []Action

	var action Action
	action.TriggerRegex = ".+asdf.+"
	action.Target = "remote"
	action.TargetParams = map[string]string{"address":"0.0.0.0:1234"}
	actions = append(actions, action)

	action.TriggerRegex = ".+abc.+"
	action.Target = "reply"
	action.TargetParams = map[string]string{"message":"{0-10}asdf"}
	actions = append(actions, action)

	connection.Actions = actions
	connections = append(connections, connection)

	connection.ConnectionType = "ssl"
	connection.Params = map[string]string{"key-path":"/tmp/key", "cert-path":"/tmp/cert"}
	connection.ProcessRegex = "."
	connection.PortRange = "100-200"

	var secondActions []Action
	action.TriggerRegex = "."
	action.Target = "remote"
	action.TargetParams = map[string]string{"address":"127.0.0.1:1758"}
	secondActions = append(secondActions, action)

	connection.Actions = secondActions
	connections = append(connections, connection)

	configTesting.Connections = connections
	SetConfigPath("/home/shahar/go/src/fyssl/config/examples/config.json")
	if !reflect.DeepEqual(configTesting, GetConfig()) {

	}






}
