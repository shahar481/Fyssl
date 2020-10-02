package utils

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/shahar481/fyssl/config"
	"math/rand"
	"net"
	"strconv"
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
		glog.Infof("Randomized address for connection %s-%s",connection.Name, connection.ListenAddress)
	}
}

func IsListening(sock net.Conn, connection *config.Connection) bool {
	if sock.LocalAddr().String() == connection.ListenAddress {
		return true
	}
	return false
}

func GetConnectionIdentifier(sock net.Conn) string {
	return sock.LocalAddr().String() + "-" + sock.RemoteAddr().String()
}

func GetActionIdentifier(sock net.Conn, action *config.Action) string {
	return GetConnectionIdentifier(sock) + "-" + strconv.Itoa(action.Id)
}