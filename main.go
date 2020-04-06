package main

import (
	"fmt"
	"fyssl/config"
)

func main() {
	config.SetConfigPath("/home/shahar/go/src/fyssl/config/examples/config.json")
	cfg := config.GetConfig()
	fmt.Printf("%+v\n", cfg)
}
