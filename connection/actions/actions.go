package actions

import (
	"fmt"
	"fyssl/config"
	"fyssl/connection/actions/targets/reply"
	"net"
	"slogger"
)

var protocolMap = map[string]func(buffer *[]byte, sock net.Conn, connection *config.Connection) {
	"remote": reply.Reply,
	"none": reply.Reply,
	"reply": reply.Reply,
}

func ProcessActions(buffer *[]byte, sock net.Conn, connection *config.Connection) {
	for index, _ := range connection.Actions {
		if connection.Actions[index].CompiledTrigger.MatchString(string(*buffer)) {
			slogger.Info(fmt.Sprintf("Catched action %s in connection %s", connection.Actions[index].Name, connection.Name))
		}
	}
}

func isListening(sock net.Conn, connection *config.Connection) bool {
	if sock.LocalAddr().String() == connection.ListenAddress {
		return true
	}
	return false
}