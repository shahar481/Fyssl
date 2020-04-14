package config

import "regexp"

// Config is the root of the config file. Meant to parse the json config file.
type Config struct {
	LogPath	string `json:"log-path"`
	Connections []Connection `json:"connections"`
}

// Connection represents a connection that fyssl handles
type Connection struct {
	Name string `json:"name"`
	ConnectionType string `json:"type"`
	RandomizeListenAddress bool
	ListenAddress string `json:"listen-address"`
	ConnectAddress string `json:"connect-address"`
	Params interface{} `json:"params"`
	Modules interface{} `json:"modules"`
	Actions []Action `json:"actions"`
}

// Action represents what we should do in each connection
type Action struct {
	Name string `json:"name"`
	TriggerRegex string `json:"trigger"`
	CompiledTrigger *regexp.Regexp
	Target string `json:"target"`
	DumpLogPackets bool `json:"dump"`
	TargetParams interface{} `json:"target-params"`
}

