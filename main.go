package main

import (
	"github.com/kaato137/quickrest/cmd"
)

var (
	Version string
	Build   string
)

func main() {
	cmd.Execute(Version, Build)
}
