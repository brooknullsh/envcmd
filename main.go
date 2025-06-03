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

func prettyPrint(commands []string, condition []string, isAsync bool) {
	var asyncFlag string
	if isAsync {
		asyncFlag = "(async) "
	}

	log.Log(log.Debug, "if \x1b[1m%s\033[0m is \x1b[1m%s\033[0m", condition[0], condition[1])
	fmt.Println("---")

	for _, cmd := range commands {
		log.Log(log.Info, "%srun \x1b[1m%s\033[0m", asyncFlag, cmd)
	}
}

func handleCommand(command string, config *config.Config) {
	switch command {
	case "create", "-c":
		config.Create()
	case "delete", "-d":
		config.Delete()
	case "show", "-s":
		for _, content := range config.Read() {
			prettyPrint(content.Commands, content.Condition, content.Async)
		}
	case "help", "-h":
		fmt.Println("Usage: envcmd COMMAND")
		fmt.Println("\nCommand line tool for running per-environment commands.")
		fmt.Println("\nOptions:")
		fmt.Println("  -h, --help  Show this message and exit.")
		fmt.Println("\nCommands:")
		fmt.Println("  -c, create  Create configuration file.")
		fmt.Println("  -d, delete  Delete configuration file.")
		fmt.Println("  -s, show    Show configuration file contents.")
	default:
		log.Abort("unknown command %s", command)
	}
}

func runAsync(content config.Content) {
	stdoutChannel := make(chan string)
	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()

		for stdout := range stdoutChannel {
			log.Log(log.Info, "%s", stdout)
		}
	}()

	for _, command := range content.Commands {
		producerWg.Add(1)
		go cmd.AsyncRun(command, stdoutChannel, &producerWg)
	}

	producerWg.Wait()
	close(stdoutChannel)
	consumerWg.Wait()
}

func run(content config.Content) {
	for _, command := range content.Commands {
		cmd.Run(command)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		log.Abort("expected 1 argument or none, got %d", len(args))
	}

	var config config.Config
	config.InitPath()

	if len(args) == 1 {
		handleCommand(args[0], &config)
		return
	}

	for _, content := range config.Read() {
		ctx, expected := content.Condition[0], content.Condition[1]

		if !context.Match(ctx, expected) {
			continue
		}

		if content.Async {
			runAsync(content)
			continue
		}

		run(content)
	}
}
