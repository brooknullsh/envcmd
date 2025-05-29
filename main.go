package main

import (
	"os"

	"github.com/brooknullsh/envcmd/internal/cmd"
	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/context"
	"github.com/brooknullsh/envcmd/internal/log"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		log.Abort("expected 1 argument or none, got %d", len(args))
	}

	var config config.Config
	config.FindPath()

	if len(args) == 1 {
		switch args[0] {
		case "create":
			config.Create()
		case "delete":
			config.Delete()
		case "show":
			for _, content := range config.Read() {
				log.PrettyContent(content.Commands, content.Condition)
			}
		default:
			log.Abort("unknown command %s", args[0])
		}

		return
	}

	for _, item := range config.Read() {
		ctx, expected := item.Condition[0], item.Condition[1]

		if !context.Match(ctx, expected) {
			log.Log(log.Warn, "%s is \x1b[1mNOT\033[0m %s", ctx, expected)
			continue
		}

		for _, command := range item.Commands {
			cmd.Run(command)
		}
	}
}
