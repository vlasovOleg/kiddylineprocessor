package kiddylineprocessor

import "time"

// Config for kiddylineprocessor
type Config struct {
	DatabaseURL string `toml:"databaseURL"`
	LogLevel    string `toml:"log_level"`

	LinesProvider struct {
		Address              string        `toml:"address"`
		RequestsTimeout      time.Duration `toml:"requestsTimeout"`
		SyncIntervalBaseball time.Duration `toml:"syncInterval_baseball"`
		SyncIntervalFootball time.Duration `toml:"syncInterval_Football"`
		SyncIntervalSoccer   time.Duration `toml:"syncInterval_Soccer"`
	} `toml:"linesProvider"`

	HTTPAPI struct {
		Address      string        `toml:"address"`
		ReadTimeout  time.Duration `toml:"readTimeout"`
		WriteTimeout time.Duration `toml:"writeTimeout"`
	} `toml:"HTTPAPI"`

	GRPC struct {
		Address string `toml:"address"`
	} `toml:"gRPC"`
}
