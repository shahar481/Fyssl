package config

// Config is the root of the config file. Meant to parse the json config file.
type Config struct {
	LogPath	string `json:"log-path"`
	Connections []Connection `json:"connections"`
}

// Connection represents a connection that fyssl handles
type Connection struct {
	ConnectionType string `json:"type"`
	ListenAddress string `json:"listen-address"`
	ConnectAddress string `json:"connect-address"`
	Params interface{} `json:"params"`
	Modules interface{} `json:"modules"`
	Actions []Action `json:"actions"`
}

// Action represents what we should do in each connection
type Action struct {
	TriggerRegex string `json:"trigger"`
	Target string `json:"target"`
	TargetParams interface{} `json:"target-params"`
}

