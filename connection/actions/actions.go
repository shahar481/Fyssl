package actions

import (
	"fmt"
	"fyssl/config"
	"fyssl/connection/actions/dump"
	"fyssl/connection/actions/targets/base"
	"fyssl/connection/actions/targets/drop"
	"fyssl/connection/actions/targets/editor"
	"fyssl/connection/actions/targets/none"
	"fyssl/connection/actions/targets/remote"
	"fyssl/connection/actions/targets/reply"
	"fyssl/connection/utils"
	"net"
	"slogger"
)

var protocolMap = map[string]func(sock net.Conn, connection *config.Connection, action *config.Action) base.Target{
	"remote": remote.CreateRemote,
	"none": none.CreateNone,
	"reply": reply.CreateReply,
	"drop": drop.CreateDrop,
	"editor": editor.CreateEditor,
}

func ProcessActions(buffer *[]byte, receiver net.Conn, sender net.Conn, activeActions *map[string]base.Target, connection *config.Connection) ([]byte, map[string]base.Target) {
	for index, _ := range connection.Actions {
		if connection.Actions[index].CompiledTrigger.MatchString(string(*buffer)) {
			slogger.Info(fmt.Sprintf("Catched action %s in connection %s", connection.Actions[index].Name, connection.Name))
			key := utils.GetActionIdentifier(receiver, &(connection.Actions[index]))
			slogger.Info(key)
			if _, ok := (*activeActions)[key]; !ok {
				(*activeActions)[key] = protocolMap[connection.Actions[index].Target](receiver, connection, &(connection.Actions[index]))
			}
			dump.Packet(buffer, receiver, sender, connection, &(connection.Actions[index]))
			buffer, err := (*activeActions)[key].ProcessTarget(buffer)
			if err != nil {
				(*activeActions)[key].Close()
				delete(*activeActions, key)
				slogger.Error(fmt.Sprintf("Error in connection: %s, in action: %s - %+v",connection.Name, connection.Actions[index].Name, err))
			}
			return *buffer, *activeActions
		}
	}
	return *buffer, *activeActions
}

