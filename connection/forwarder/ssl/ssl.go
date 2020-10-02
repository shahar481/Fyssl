package ssl

import (
	"crypto/tls"
	"fmt"
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection/forwarder"
	"github.com/shahar481/fyssl/connection/forwarder/base"
	"github.com/shahar481/fyssl/connection/utils"
	"net"
	"slogger"
	"time"
)

const (
	protocol = "tcp"
)

type Ssl struct {
	Connection *config.Connection
	listener net.Listener
	tlsCfg *tls.Config
}

func NewSslForwarder(connection *config.Connection) base.Listener {
	s := Ssl{
		Connection: connection,
	}
	return s
}

func (s Ssl) StartListening() {
	slogger.Info(fmt.Sprintf("Initializing an ssl connection-%s",s.Connection.Name))
	s.createTlsConfig()
	s.startListeningSocket()
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			slogger.Error(fmt.Sprintf("Error accepting incoming connection at %s-%+v", s.Connection.Name, err))
			continue
		}
		go s.forwardConnection(conn)
	}
}

func (s *Ssl) startListeningSocket() {
	for {
		l, err := tls.Listen(protocol, s.Connection.ListenAddress, s.tlsCfg)
		if err != nil {
			slogger.Error(fmt.Sprintf("Error occured on %s on address %s-%+v", s.Connection.Name, s.Connection.ListenAddress,err))
			time.Sleep(forwarder.ListenErrorTimeout * time.Second)
			utils.SetListenAddress(s.Connection)
		} else {
			slogger.Info(fmt.Sprintf("Started listening on connection %s on address %s", s.Connection.Name, s.Connection.ListenAddress))
			s.listener = l
			return
		}
	}
}

func (s *Ssl) createTlsConfig() {
	tlsParams := s.Connection.Params.(map[string]interface{})
	for {
		cer, err := tls.LoadX509KeyPair(tlsParams["cert-path"].(string), tlsParams["key-path"].(string))
		if err != nil {
			slogger.Error(fmt.Sprintf("Error occured loading tls keys at connection:%s-%+v",s.Connection.Name, err))
			time.Sleep(forwarder.ListenErrorTimeout * time.Second)
			continue
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}, InsecureSkipVerify: true}
		s.tlsCfg = config
		return
	}

}

func (s Ssl) forwardConnection(receiver net.Conn) {
	sender, err := tls.Dial(protocol, s.Connection.ConnectAddress, s.tlsCfg)
	if err != nil {
		slogger.Error(fmt.Sprintf("Couldn't connect to %s on connection %s-%+v", s.Connection.ConnectAddress, s.Connection.Name, err))
		receiver.Close()
		return
	}
	forwarder.StartForwardSockets(receiver, sender, s.Connection)
}