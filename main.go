package main

import (
	"fmt"
	"os"

	"github.com/brooknullsh/envcmd/internal/cmd"
	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/log"
)

func start(cfg *config.Config) {
	for _, obj := range cfg.Read() {
		cmd.Run(obj)
	}
}

func main() {
	var cfg config.Config
	cfg.Init()

	args := os.Args[1:]
	if len(args) == 0 {
		start(&cfg)
		return
	}

	switch args[0] {
	case "create", "-c":
		cfg.Create()
	case "delete", "-d":
		cfg.Delete()
	case "list", "-l":
		cfg.List()
	case "help", "-h":
		fmt.Println("Usage: envcmd [COMMAND]")
		fmt.Println("\nCommand line tool for running per-environment commands.")
		fmt.Println("\nCommands:")
		fmt.Println("  create, -c  Create configuration file.")
		fmt.Println("  delete, -d  Delete configuration file.")
		fmt.Println("  list,   -l  Show configuration file contents.")
		fmt.Println("  help,   -h  Show this message.")
	default:
		log.Warn("unknown command -> %s", args[0])
	}
}
