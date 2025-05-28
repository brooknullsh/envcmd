package cmd

import (
	"bufio"
	"os"
	"os/exec"

	"github.com/brooknullsh/envcmd/internal/log"
)

func extractOutput(scanner *bufio.Scanner) {
	for scanner.Scan() {
		log.Log(log.Info, "%s", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Log(log.Warn, "reading stdout: %v", err)
	}
}

func Run(command string) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Log(log.Warn, "failed to pipe stdout: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		log.Log(log.Error, "starting command: %v", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(stdout)
	extractOutput(scanner)

	err = cmd.Wait()
	if err != nil {
		log.Log(log.Error, "exited: %v", err)
		os.Exit(1)
	}
}
