package main

import (
	"github.com/faireal/kobe/cmd/client/root"
	"github.com/faireal/kobe/pkg/config"
	"os"
)

func main() {
	config.Init()
	if err := root.Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
