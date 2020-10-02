package remote

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection/actions/targets/base"
	"github.com/shahar481/fyssl/connection/utils"
	"net"
	"slogger"
)

const (
	remoteConnectionType = "tcp"
	listenerAddressKey = "listener-socket-address"
	senderAddressKey = "sender-socket-address"
	socketListKey = "sockets"
	lengthSize = 2
)

type Remote struct {
	sock net.Conn
	remoteSock net.Conn
	connection *config.Connection
	action *config.Action
	remoteAddress string
}

func (r *Remote) Close() {
	if r.remoteSock != nil {
		r.remoteSock.Close()
	}
}

func CreateRemote(sock net.Conn, connection *config.Connection, action *config.Action) base.Target {
	var target base.Target
	var addr string
	dict, _ := action.TargetParams.(map[string]interface{})
	if utils.IsListening(sock, connection) {
		addr = dict[listenerAddressKey].(string)
	} else {
		addr = dict[senderAddressKey].(string)
	}
	target = &Remote{
		sock:          sock,
		connection:    connection,
		action:        action,
		remoteAddress: addr,
	}
	return target
}
func (r *Remote) ProcessTarget(buffer *[]byte) (*[]byte, error) {
	return r.startRemoteForwarding(buffer)
}

func (r *Remote) printError(err error) {
	slogger.Info(fmt.Sprintf("Error occured in connection: %s, action: %s - %+v", r.connection.Name, r.action.Name, err))
}

func (r *Remote) startRemoteForwarding(buffer *[]byte) (*[]byte, error) {
	if r.remoteSock == nil {
		err := r.createRemoteConnection()
		if err != nil {
			return buffer, err
		}
	}
	finalBuffer, err := r.forwardToRemoteConnection(buffer)
	return finalBuffer, err
}

func (r *Remote) forwardToRemoteConnection(buffer *[]byte) (*[]byte, error) {
	var length = make([]byte, lengthSize)
	binary.BigEndian.PutUint16(length, uint16(len(*buffer)))
	_, err := r.remoteSock.Write(append(length,*buffer ...))
	if err != nil {
		return buffer, err
	}
	buff, err := r.getReceivedMessage()
	if err != nil {
		return buffer, err
	}
	return buff, err
}

func (r *Remote) getReceivedMessage() (*[]byte, error) {
	var buff = make([]byte, 0)
	expectedLength, err := r.getMessageLength()
	if err != nil {
		return &buff, err
	}
	received := 0
	for received < expectedLength {
		var receivedBuff = make([]byte, 1024)
		len, err := r.remoteSock.Read(receivedBuff)
		if err != nil {
			return &buff, err
		}
		buff = append(buff, receivedBuff[:len] ...)
		received += len
	}
	return &buff, err
}

func (r *Remote) getMessageLength() (int, error) {
	var length = make([]byte, lengthSize)
	gotLength, err := r.remoteSock.Read(length)
	if err != nil {
		return 0, err
	}
	if gotLength < lengthSize {
		return 0, errors.New("got length less than 2 bytes in received message")
	}
	return int(binary.BigEndian.Uint16(length)), err
}



func (r *Remote) createRemoteConnection() error {
	conn, err := net.Dial(remoteConnectionType, r.remoteAddress)
	if err != nil {
		return err
	}
	r.remoteSock = conn
	return nil
}
