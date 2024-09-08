package config

type ServerConfig struct {
	ServerAddr    string `env:"SERVER_ADDRESS"`
	TelemetryAddr string `env:"TELEMETRY_ADDRESS"`
}