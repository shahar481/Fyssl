package utils

import (
	"fmt"
	"fyssl/config"
	"math/rand"
	"slogger"
	"time"
)

const (
	listenIp = "127.0.0.1"
	maxPort = 65535
	minPort = 1
)

func randomizeListenAddress() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s:%d", listenIp, rand.Intn(maxPort - minPort) + minPort)
}

func SetListenAddress(connection *config.Connection) {
	if connection.RandomizeListenAddress {
		connection.ListenAddress = randomizeListenAddress()
		slogger.Debug(fmt.Sprintf("Randomized address for connection %s-%s",connection.Name, connection.ListenAddress))
	}
}