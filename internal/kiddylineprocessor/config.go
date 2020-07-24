package kiddylineprocessor

import "time"

// Config for kiddylineprocessor
type Config struct {
	DatabaseURL              string        `toml:"databaseURL"`
	UpdateByProviderBaseball time.Duration `toml:"updateByProvider_baseball"`
	UpdateByProviderFootball time.Duration `toml:"updateByProvider_football"`
	UpdateByProviderSoccer   time.Duration `toml:"updateByProvider_soccer"`
	HTTPserverIP             string        `toml:"HTTPserver_ip"`
	GRPCserverIP             string        `toml:"GRPCserver_ip"`
	LogLevel                 string        `toml:"log_level"`
}
