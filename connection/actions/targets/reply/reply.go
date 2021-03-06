package reply

import (
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection/actions/targets/base"
	"github.com/shahar481/fyssl/connection/actions/targets/reply/language"
	"net"
)

type Reply struct {
	connection *config.Connection
	action *config.Action
}

func (r *Reply) ProcessTarget(buffer *[]byte) (*[]byte, error) {
	dict, _ := r.action.TargetParams.(map[string]interface{})
	return language.Process(buffer, dict["message"].(string))
}

func (r *Reply) Close() {}

func CreateReply(sock net.Conn, connection *config.Connection, action *config.Action) base.Target {
	r := Reply{
		connection:connection,
		action:action,
	}
	return &r
}