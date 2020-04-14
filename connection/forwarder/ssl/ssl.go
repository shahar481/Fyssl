package ssl

import (
	"crypto/tls"
	"fmt"
	"fyssl/config"
	"fyssl/connection/forwarder"
	"fyssl/connection/utils"
	"net"
	"slogger"
	"time"
)

const (
	protocol = "tcp"
)

func StartListening(connection *config.Connection) {
	slogger.Info(fmt.Sprintf("Initializing an ssl connection-%s",connection.Name))
	tlsCfg := createTlsConfig(connection)
	l := getListeningSocket(connection, tlsCfg)
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			slogger.Error(fmt.Sprintf("Error accepting incoming connection at %s-%+v", connection.Name, err))
			continue
		}
		go forwardConnection(conn, connection, tlsCfg)
	}
}

func getListeningSocket(connection *config.Connection, tlsCfg *tls.Config) net.Listener {
	for {
		l, err := tls.Listen(protocol, connection.ListenAddress, tlsCfg)
		if err != nil {
			slogger.Error(fmt.Sprintf("Error occured on %s on address %s-%+v",connection.Name,connection.ListenAddress,err))
			time.Sleep(forwarder.ListenErrorTimeout * time.Second)
			utils.SetListenAddress(connection)
		} else {
			slogger.Info(fmt.Sprintf("Started listening on connection %s on address %s", connection.Name, connection.ListenAddress))
			return l
		}
	}
}

func createTlsConfig(connection *config.Connection) *tls.Config {
	tlsParams := connection.Params.(map[string]interface{})
	for {
		cer, err := tls.LoadX509KeyPair(tlsParams["cert-path"].(string), tlsParams["key-path"].(string))
		if err != nil {
			slogger.Error(fmt.Sprintf("Error occured loading tls keys at connection:%s-%+v",connection.Name, err))
			time.Sleep(forwarder.ListenErrorTimeout * time.Second)
			continue
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}, InsecureSkipVerify: true}
		return config
	}

}

func forwardConnection(receiver net.Conn, connection *config.Connection, tlsCfg *tls.Config) {
	sender, err := tls.Dial(protocol, connection.ConnectAddress, tlsCfg)
	if err != nil {
		slogger.Error(fmt.Sprintf("Couldn't connect to %s on connection %s-%+v", connection.ConnectAddress, connection.Name, err))
		receiver.Close()
		return
	}
	forwarder.StartForwardSockets(receiver, sender, connection)
}