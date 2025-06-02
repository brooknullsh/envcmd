package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/brooknullsh/envcmd/internal/cmd"
	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/context"
	"github.com/brooknullsh/envcmd/internal/log"
)

func handleCommand(cmdStr string, c *config.Config) {
	switch cmdStr {
	case "create", "-c":
		c.Create()
	case "delete", "-d":
		c.Delete()
	case "show", "-s":
		for _, cont := range c.Read() {
			log.PrettyContent(cont.Commands, cont.Condition)
		}
	case "-h", "--help":
		fmt.Println("Usage: envcmd COMMAND")
		fmt.Println("\nCommand line tool for running per-environment commands.")
		fmt.Println("\nOptions:")
		fmt.Println("  -h, --help  Show this message and exit.")
		fmt.Println("\nCommands:")
		fmt.Println("  -c, create  Create configuration file.")
		fmt.Println("  -d, delete  Delete configuration file.")
		fmt.Println("  -s, show    Show configuration file contents.")
	default:
		log.Abort("unknown command %s", cmdStr)
	}
}

func runAsync(cont config.Content) {
	ch := make(chan string)
	var prodWg sync.WaitGroup
	var consWg sync.WaitGroup

	consWg.Add(1)
	go func() {
		defer consWg.Done()

		for out := range ch {
			log.Log(log.Info, "%s", out)
		}
	}()

	for idx, cmdStr := range cont.Commands {
		prodWg.Add(1)
		go cmd.Run(idx, cmdStr, ch, &prodWg)
	}

	prodWg.Wait()
	close(ch)
	consWg.Wait()
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		log.Abort("expected 1 argument or none, got %d", len(args))
	}

	var c config.Config
	c.InitPath()

	if len(args) == 1 {
		handleCommand(args[0], &c)
		return
	}

	for _, cont := range c.Read() {
		ctx, exp := cont.Condition[0], cont.Condition[1]

		if !context.Match(ctx, exp) {
			log.Log(log.Warn, "%s is \x1b[1mNOT\033[0m %s", ctx, exp)
			continue
		}

		runAsync(cont)
	}
}
