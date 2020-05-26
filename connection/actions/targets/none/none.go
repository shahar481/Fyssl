package none

import (
	"fyssl/config"
	"fyssl/connection/actions/targets/base"
	"net"
)

type None struct {}

func (n *None) ProcessTarget(buffer *[]byte) (*[]byte, error) {
	return buffer, nil
}

func (n *None) Close() {}

func CreateNone(sock net.Conn, connection *config.Connection, action *config.Action) base.Target {
	n := None{}
	return &n
}