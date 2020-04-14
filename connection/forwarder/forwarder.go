package forwarder

import (
	"fmt"
	"fyssl/config"
	"fyssl/connection/actions"
	"net"
	"slogger"
)

func StartForwardSockets(first net.Conn, second net.Conn, connection *config.Connection) {
	slogger.Info(fmt.Sprintf("Forwarding a connection in %s", connection.Name))
	go forwardSockets(first, second, connection)
	go forwardSockets(second, first, connection)
}

func forwardSockets(receiver net.Conn, sender net.Conn, connection *config.Connection) {
	var buf = make([]byte, 1024)
	for {
		copy(buf, make([]byte, len(buf)))
		_, err := receiver.Read(buf)
		if err != nil {
			slogger.Info(fmt.Sprintf("Connection closed on %s", connection.Name))
			receiver.Close()
			sender.Close()
			return
		}
		actions.ProcessActions(&buf, receiver, connection)
		_, err = sender.Write(buf)
		if err != nil {
			receiver.Close()
			sender.Close()
			return
		}
	}
}