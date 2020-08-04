package main

import (
	"flag"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/kiddylineprocessor"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "path to configs folder")
	flag.Parse()

	config := kiddylineprocessor.Config{}
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Panic(err)
	}
	config.LinesProvider.RequestsTimeout *= time.Second
	config.LinesProvider.SyncIntervalBaseball *= time.Second
	config.LinesProvider.SyncIntervalFootball *= time.Second
	config.LinesProvider.SyncIntervalSoccer *= time.Second
	config.HTTPAPI.ReadTimeout *= time.Second
	config.HTTPAPI.WriteTimeout *= time.Second

	klp := kiddylineprocessor.New(&config)
	klp.Start()
}
