package connection

import (
	"bufio"
	"fmt"
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection/forwarder/base"
	"github.com/shahar481/fyssl/connection/forwarder/ssl"
	"github.com/shahar481/fyssl/connection/forwarder/tcp"
	"os"
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
		fmt.Printf(reader.ReadString('\n'))
	}
}

