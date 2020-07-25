package kiddylineprocessor

import "time"

// Config for kiddylineprocessor
type Config struct {
	DatabaseURL                  string        `toml:"databaseURL"`
	LinesProviderAddress         string        `toml:"linesProvider_address"`
	LinesProviderRequestsTimeout time.Duration `toml:"linesProvider_requestsTimeout"`
	UpdateByProviderBaseball     time.Duration `toml:"linesProvider_baseball"`
	UpdateByProviderFootball     time.Duration `toml:"linesProvider_football"`
	UpdateByProviderSoccer       time.Duration `toml:"linesProvider_soccer"`
	HTTPserverIP                 string        `toml:"HTTPserver_ip"`
	GRPCserverIP                 string        `toml:"GRPCserver_ip"`
	LogLevel                     string        `toml:"log_level"`
}
