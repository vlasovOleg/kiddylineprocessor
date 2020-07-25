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

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "path to configs folder")
}

func main() {
	flag.Parse()

	config := kiddylineprocessor.Config{}
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	klp := kiddylineprocessor.New(&config)
	klp.Start()
	for {
	}
}
