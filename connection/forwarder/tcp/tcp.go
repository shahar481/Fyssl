package tcp

import (
	"fmt"
	"fyssl/config"
	"fyssl/connection/forwarder"
	"fyssl/connection/utils"
	"net"
	"slogger"
	"time"
)

const (
	connType = "tcp"
)

func StartListening(connection *config.Connection) {
	slogger.Info(fmt.Sprintf("Initializing a tcp connection-%s", connection.Name))
	l := getListeningSocket(connection)
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			slogger.Error(fmt.Sprintf("Error accepting incomming connection at %s-%+v",connection.Name, err))
		}
		go forwardConnection(conn, connection)
	}
}

func getListeningSocket(connection *config.Connection) net.Listener {
	for {
		l, err := net.Listen(connType, connection.ListenAddress)
		if err != nil {
			slogger.Error(fmt.Sprintf("Error occured on connection %s on address %s-%+v",connection.Name, connection.ListenAddress, err))
			time.Sleep(forwarder.ListenErrorTimeout * time.Second)
			utils.SetListenAddress(connection)
		} else {
			slogger.Info(fmt.Sprintf("Started listening on connection %s on address %s", connection.Name, connection.ListenAddress))
			return l
		}
	}
}

func forwardConnection(receiver net.Conn, connection *config.Connection) {
	sender, err := net.Dial(connType, connection.ConnectAddress)
	if err != nil {
		slogger.Error(fmt.Sprintf("Couldn't connect to %s on connection %s-%+v", connection.ConnectAddress, connection.Name, err))
		receiver.Close()
		return
	}
	forwarder.StartForwardSockets(receiver, sender, connection)
}