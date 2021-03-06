package tcp

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection/forwarder"
	"github.com/shahar481/fyssl/connection/forwarder/base"
	"github.com/shahar481/fyssl/connection/utils"
	"log"
	"net"
	"time"
)

const (
	connType = "tcp"
)


type Tcp struct {
	Connection *config.Connection
	listener net.Listener
}

func NewTcpForwarder(connection *config.Connection) base.Listener {
	t := Tcp{
		Connection: connection,
	}
	return t
}

func (t Tcp) StartListening() {
	log.Printf(fmt.Sprintf("Initializing a tcp connection-%s", t.Connection.Name))
	t.startListeningSocket()
	defer t.listener.Close()
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			glog.Infof("Error accepting incomming connection at %s-%+v",t.Connection.Name, err)
		}
		go t.forwardConnection(conn)
	}
}

func (t *Tcp) startListeningSocket() {
	for {
		l, err := net.Listen(connType, t.Connection.ListenAddress)
		if err != nil {
			glog.Errorf("Error occured on connection %s on address %s-%+v",t.Connection.Name, t.Connection.ListenAddress, err)
			time.Sleep(forwarder.ListenErrorTimeout * time.Second)
			utils.SetListenAddress(t.Connection)
		} else {
			glog.Errorf("Started listening on connection %s on address %s", t.Connection.Name, t.Connection.ListenAddress)
			t.listener = l
			return
		}
	}
}

func (t Tcp) forwardConnection(receiver net.Conn) {
	sender, err := net.Dial(connType, t.Connection.ConnectAddress)
	if err != nil {
		glog.Errorf("Couldn't connect to %s on connection %s-%+v", t.Connection.ConnectAddress, t.Connection.Name, err)
		receiver.Close()
		return
	}
	forwarder.StartForwardSockets(receiver, sender, t.Connection)
}