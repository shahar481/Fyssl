package main

import (
	"fyssl/config"
	"log"
)

func main() {
	config.SetConfigPath("/home/shahar/go/src/fyssl/examples/config.json")
	cfg := config.GetConfig()
	log.Println(cfg)
}
