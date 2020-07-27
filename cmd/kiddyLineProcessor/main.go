package main

import (
	"flag"
	"log"

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

	klp := kiddylineprocessor.New(&config)
	klp.Start()
}
