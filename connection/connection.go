package connection

import (
	"bufio"
	"fyssl/config"
	"fyssl/connection/forwarder/base"
	"fyssl/connection/forwarder/ssl"
	"fyssl/connection/forwarder/tcp"
	"os"
	"slogger"
)

var protocolMap = map[string]func(config *config.Connection) base.Listener {
	"tcp": tcp.NewTcpForwarder,
	"ssl": ssl.NewSslForwarder,
}

func StartConnections() {
	activeConnections := make(map[int]base.Listener)
	cfg := config.GetConfig()
	for index, _ := range config.GetConfig().Connections {
		activeConnections[cfg.Connections[index].Id] = protocolMap[cfg.Connections[index].ConnectionType](&cfg.Connections[index])
		go activeConnections[cfg.Connections[index].Id].StartListening()
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		slogger.Info(reader.ReadString('\n'))
	}
}

