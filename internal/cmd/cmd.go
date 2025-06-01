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

func Run(index int, command string, channel chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Log(log.Warn, "failed to pipe stdout: %v", err)
	}

	if err = cmd.Start(); err != nil {
		log.Abort("starting command: %v", err)
	}

	colour := nextColour()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		channel <- fmt.Sprintf("%s(%d)\033[0m %s", colour, index, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Log(log.Warn, "reading stdout: %v", err)
	}

	if err = cmd.Wait(); err != nil {
		log.Abort("exited: %v", err)
	}
}
