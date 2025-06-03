package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/brooknullsh/envcmd/internal/log"
)

var (
	colourIndex int
	colourMutex sync.Mutex
)

func nextColour() string {
	colourMutex.Lock()
	defer colourMutex.Unlock()

	colour := log.Colours[colourIndex]
	colourIndex = (colourIndex + 1) % len(log.Colours)
	return colour
}

func sharedRun(command string, handleStdout func(stdout string)) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Log(log.Warn, "failed to pipe stdout: %v", err)
	}

	if err = cmd.Start(); err != nil {
		log.Abort("starting command: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		handleStdout(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Log(log.Warn, "reading stdout: %v", err)
	}

	if err = cmd.Wait(); err != nil {
		log.Abort("exited: %v", err)
	}
}

func AsyncRun(command string, stdoutChannel chan<- string, producerWg *sync.WaitGroup) {
	defer producerWg.Done()

	colour := nextColour()
	log.Log(log.Debug, "running %s%s...", colour, command)

	handleStdout := func(stdout string) {
		stdoutChannel <- fmt.Sprintf("%s%s", colour, stdout)
	}

	sharedRun(command, handleStdout)
}

func Run(command string) {
	log.Log(log.Debug, "running %s...", command)

	handleStdout := func(stdout string) {
		log.Log(log.Info, "%s", stdout)
	}

	sharedRun(command, handleStdout)
}
