package main

import (
	"os"

	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/log"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Log(log.Error, "Expected 1 argument, got %d", len(args))
		return
	}

	log.Log(log.Info, "%s", args[0])
	config.Run()
}
