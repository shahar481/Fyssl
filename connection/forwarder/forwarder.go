package forwarder

import (
	"fmt"
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection/actions"
	"github.com/shahar481/fyssl/connection/actions/targets/base"
	"net"
	"slogger"
)

func StartForwardSockets(first net.Conn, second net.Conn, connection *config.Connection) {
	slogger.Info(fmt.Sprintf("Forwarding a connection in %s", connection.Name))
	go forwardSockets(first, second, connection)
	go forwardSockets(second, first, connection)
}

func forwardSockets(receiver net.Conn, sender net.Conn, connection *config.Connection) {
	ActiveActions := make(map[string]base.Target)
	for {
		var buf = make([]byte, 1024)
		copy(buf, make([]byte, len(buf)))
		length, err := receiver.Read(buf)
		cutMessage :=  buf[:length]
		buf = nil
		if err != nil {
			closeConnection(receiver, sender, ActiveActions)
			return
		}
		cutMessage, ActiveActions = actions.ProcessActions(&cutMessage, receiver, sender, &ActiveActions, connection)
		_, err = sender.Write(cutMessage)
		if err != nil {
			closeConnection(receiver, sender, ActiveActions)
			return
		}
	}
}

func closeConnection(receiver net.Conn, sender net.Conn, ActiveActions map[string]base.Target) {
	receiver.Close()
	sender.Close()
	for _, action := range ActiveActions {
		action.Close()
	}
}