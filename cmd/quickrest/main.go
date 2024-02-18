package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/kaato137/quickrest/internal"
)

var (
	Version string
	Build   string
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "quickrest.yml", "path to a quick rest config")
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

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

func printVersion() {
	fmt.Printf("QuickREST v%s.%s", Version, Build)
}
