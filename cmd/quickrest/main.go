package main

import (
	"flag"
	"net/http"

	"github.com/kaato137/quickrest/internal"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "quickrest.yml", "path to a quick rest config")
	flag.Parse()

	cfg, err := internal.LoadConfigFromFile(configPath)
	if err != nil {
		panic(err)
	}

	server, err := internal.NewServerFromConfig(cfg)
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(cfg.Address, server); err != nil {
		panic(err)
	}
}
