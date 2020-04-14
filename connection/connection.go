package connection

import (
	"bufio"
	"fyssl/config"
	"fyssl/connection/forwarder/ssl"
	"fyssl/connection/forwarder/tcp"
	"os"
	"slogger"
)

var protocolMap = map[string]func(config *config.Connection){
	"tcp": tcp.StartListening,
	"ssl": ssl.StartListening,
}

func StartConnections() {
	cfg := config.GetConfig()
	for index, _ := range config.GetConfig().Connections {
		go protocolMap[cfg.Connections[index].ConnectionType](&cfg.Connections[index])
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		slogger.Info(reader.ReadString('\n'))
	}
}

