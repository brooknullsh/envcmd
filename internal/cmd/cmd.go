package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/brooknullsh/envcmd/internal/log"
)

var (
	colIdx int
	colMtx sync.Mutex
)

func nextColour() string {
	colMtx.Lock()
	defer colMtx.Unlock()

	col := log.Colours[colIdx]
	colIdx = (colIdx + 1) % len(log.Colours)
	return col
}

func Run(idx int, cmdStr string, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("bash", "-c", cmdStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Log(log.Warn, "failed to pipe stdout: %v", err)
	}

	if err = cmd.Start(); err != nil {
		log.Abort("starting command: %v", err)
	}

	col := nextColour()
	s := bufio.NewScanner(stdout)
	for s.Scan() {
		ch <- fmt.Sprintf("%s(%d)\033[0m %s", col, idx, s.Text())
	}

	if err := s.Err(); err != nil {
		log.Log(log.Warn, "reading stdout: %v", err)
	}

	if err = cmd.Wait(); err != nil {
		log.Abort("exited: %v", err)
	}
}
