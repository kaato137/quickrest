package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kaato137/quickrest/internal"
	"github.com/kaato137/quickrest/internal/conf"
)

var (
	Version string
	Build   string
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "path to a quick rest config")
	flag.Parse()

	shouldClose, exitCode := processCmd()
	if shouldClose {
		os.Exit(exitCode)
	}

	cfg, err := conf.LoadConfigFromFile(configPath)
	if err != nil {
		if errors.Is(err, conf.ErrDefaultPathNotFound) {
			printMessageAboutDefaultConfig()
			os.Exit(1)
		}
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

func processCmd() (shouldClose bool, exitCode int) {
	cmd := strings.ToLower(flag.Arg(0))
	switch cmd {

	case "version":
		printVersion()
		return true, 0

	case "gdc", "generate-default-config":
		if err := conf.GenerateDefault(); err != nil {
			if errors.Is(err, conf.ErrDefaultConfigAlreadyExists) {
				fmt.Fprintln(os.Stderr, "[ERR] Default config already exists.")
				return true, 1
			}
			panic(err)
		}
		return true, 0
	}

	return false, 0
}

func printVersion() {
	fmt.Printf("QuickREST v%s.%s\n", Version, Build)
}

func printMessageAboutDefaultConfig() {
	fmt.Fprintln(os.Stderr, `[ERR] Configuration file not found.
The expected default configuration file is named 'quickrest.yml' or 'quickrest.yaml'.
Alternatively, you can specify a custom configuration file using:
	quickrest -c your_custom_config.yml

You can also generate the default configuration file by running:
	quickrest generate-default-config
or
	quickrest gdc
	`)
}
