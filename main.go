package main

import (
	"fmt"
	"os"

	"github.com/brooknullsh/envcmd/internal/cmd"
	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/log"
)

func main() {
	var config config.Config
	config.InitPath()

	args := os.Args[1:]
	if len(args) == 0 {
		for _, content := range config.Read() {
			cmd.Run(content)
		}
		return
	}

	switch args[0] {
	case "create", "-c":
		config.Create()
	case "delete", "-d":
		config.Delete()
	case "list", "-l":
		for _, content := range config.Read() {
			content.Print()
			fmt.Println("---")
		}
	case "help", "-h":
		fmt.Println("Usage: envcmd COMMAND")
		fmt.Println("\nCommand line tool for running per-environment commands.")
		fmt.Println("\nCommands:")
		fmt.Println("  help,   -h  Show this message and exit.")
		fmt.Println("  create, -c  Create configuration file.")
		fmt.Println("  delete, -d  Delete configuration file.")
		fmt.Println("  show,   -s  Show configuration file contents.")
	default:
		log.Warn("unknown command -> %s", args[0])
	}
}
