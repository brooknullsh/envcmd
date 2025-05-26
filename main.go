package main

import (
	"os"

	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/context"
	"github.com/brooknullsh/envcmd/internal/log"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		log.Log(log.Error, "expected 1 argument or none, got %d", len(args))
		os.Exit(1)
	}

	var config config.Config
	config.FindPath()

	if len(args) == 0 {
		content := config.Read()

		for _, item := range content {
			ctx, expected := item.Condition[0], item.Condition[1]

			if context.Match(ctx, expected) {
				item.Process(true)
			} else {
				log.Log(log.Warn, "%s is \x1b[1mNOT\033[0m %s", ctx, expected)
			}
		}

		return
	}

	switch args[0] {
	case "create":
		config.Create()
	case "delete":
		config.Delete()
	case "show":
		content := config.Read()
		for _, item := range content {
			item.Process(false)
		}
	default:
		log.Log(log.Error, "unknown command %s", args[0])
		os.Exit(1)
	}
}
