package main

import (
	"appliedgo.net/what"
	"fmt"
	"github.com/shahar481/fyssl/config"
	"github.com/shahar481/fyssl/connection"
	"log"
	"os"
)

func processArgs() {
	processHelp()
	processConfig()
}

func printHelp() {
	fmt.Println("-h   Help" +
		"\n-c   ConfigPath, usage: -c {path}")
}

func processHelp() {
	if doesArgExist("-h") != -1 {
		printHelp()
		log.Fatal("No config file found")
	}
}

func processConfig() {
	if doesArgExist("-c") != -1 {
		if doesArgExist("-c") + 1 <= len(os.Args[1:]) - 1{
			configPath := os.Args[1:][doesArgExist("-c") + 1]
			config.SetConfigPath(configPath)
			slogger.Info(fmt.Sprintf("Set config path to:%s", configPath))
		} else {
			printHelp()
			log.Fatal("No config file found")
		}
	} else {
		printHelp()
		log.Fatal("No config file found")
	}
}

func doesArgExist(arg string) int {
	args := os.Args[1:]
	for index, val := range args {
		if arg == val {
			return index
		}
	}
	return -1
}

func main() {
	processArgs()
	connection.StartConnections()
}


