package cmd

import (
	"bufio"
	"fmt"
	"io"
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

func sharedScanner(stream io.ReadCloser, fn func(out string), wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		fn(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("reading stream: %v\n", err)
	}
}

func sharedRun(command string, fn func(out string)) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Warn("piping stdout -> %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Warn("piping stderr -> %v", err)
	}

	if err = cmd.Start(); err != nil {
		log.Abort("starting command -> %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go sharedScanner(stdout, fn, &wg)

	wg.Add(1)
	go sharedScanner(stderr, fn, &wg)

	wg.Wait()
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
		log.Warn("no match for \x1b[1m%s\033[0m", content.Name)
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
