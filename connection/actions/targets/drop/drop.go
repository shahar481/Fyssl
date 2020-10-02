package drop

import (
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection/actions/targets/base"
	"net"
)

type Drop struct {}

func (d *Drop) ProcessTarget(buffer *[]byte) (*[]byte, error) {
	buff := make([]byte, 0)
	return &buff, nil
}

func (d *Drop) Close() {}

func CreateDrop(sock net.Conn, connection *config.Connection, action *config.Action) base.Target {
	d := Drop{}
	return &d
}