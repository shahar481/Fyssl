package main

import (
	"fyssl/config"
	"fyssl/connection"
)

func main() {
	config.SetConfigPath("/home/shahar/go/src/fyssl/config/examples/config.json")
	connection.StartConnections()
}
