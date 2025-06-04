package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/context"
	"github.com/brooknullsh/envcmd/internal/log"
)

var (
	colourIndex int
	colourMutex sync.Mutex
)

func asyncColour() string {
	colourMutex.Lock()
	defer colourMutex.Unlock()

	colour := log.Colours[colourIndex]
	colourIndex = (colourIndex + 1) % len(log.Colours)
	return colour
}

func sharedRun(command string, fn func(out string)) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Warn("failed to pipe stdout -> %v", err)
	}

	if err = cmd.Start(); err != nil {
		log.Abort("starting command -> %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fn(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Warn("reading stdout -> %v", err)
	}

	if err = cmd.Wait(); err != nil {
		log.Abort("exited -> %v", err)
	}
}

func syncRun(content config.Content) {
	fn := func(out string) {
		fmt.Fprintf(os.Stdout, "%s\n", out)
	}

	for _, cmd := range content.Commands {
		sharedRun(cmd, fn)
	}
}

func asyncRun(command string, stdout chan<- string, producerWg *sync.WaitGroup) {
	defer producerWg.Done()

	colour := asyncColour()
	fn := func(out string) {
		stdout <- fmt.Sprintf("%s%s\033[0m", colour, out)
	}

	sharedRun(command, fn)
}

func readStdout(stdout <-chan string, consumerWg *sync.WaitGroup) {
	defer consumerWg.Done()

	for out := range stdout {
		fmt.Fprintf(os.Stdout, "%s\n", out)
	}
}

func Run(content config.Content) {
	if !context.Match(content) {
		log.Warn("no \x1b[1m%s\033[0m match for \x1b[1m%s\033[0m", content.Context, content.Target)
		return
	}

	log.Info("matched with \x1b[1m%s\033[0m", content.Name)

	if !content.Async {
		syncRun(content)
		return
	}

	stdout := make(chan string)
	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	consumerWg.Add(1)
	go readStdout(stdout, &consumerWg)

	for _, cmd := range content.Commands {
		producerWg.Add(1)
		go asyncRun(cmd, stdout, &producerWg)
	}

	producerWg.Wait()
	close(stdout)
	consumerWg.Wait()
}
