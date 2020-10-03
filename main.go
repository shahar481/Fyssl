package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection"
)

func main() {
	processArgs()
	connection.StartConnections()
}

func processArgs() {
	configPath := flag.String("c", "", "ConfigPath, usage: -c {path}")
	help := flag.Bool("h", false, "Help")
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *configPath != "" {
		config.SetConfigPath(*configPath)
		glog.Infof("Set config path to:%s", *configPath)
		return
	}

	printHelp()
	glog.Fatal("No config file found")
}

func printHelp() {
	fmt.Println("-h   Help")
	fmt.Println("-c   ConfigPath, usage: -c {path}")
}
