package main

import (
	"os"

	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/log"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Log(log.Error, "expected 1 argument, got %d", len(args))
		os.Exit(1)
	}

	var config config.Config
	config.FindPath()

	switch args[0] {
	case "create":
		config.Create()
	case "delete":
		config.Delete()
	case "list":
		config.List()
	default:
		log.Log(log.Error, "unknown command %s", args[0])
		os.Exit(1)
	}
}
