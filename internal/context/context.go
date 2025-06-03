package context

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/brooknullsh/envcmd/internal/log"
)

func branchMatch(actual string) bool {
	command := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	stdout, err := command.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 128 {
			log.Log(log.Debug, "no git in current directory")
			return false
		}

		log.Abort("reading git branch: %v", err)
	}

	return strings.TrimSpace(string(stdout)) == actual
}

func directoryMatch(actual string) bool {
	dirPath, err := os.Getwd()
	if err != nil {
		log.Abort("reading current directory: %v", err)
	}

	return filepath.Base(dirPath) == actual
}

func Match(context string, actual string) bool {
	if context == "directory" {
		return directoryMatch(actual)
	} else if context == "branch" {
		return branchMatch(actual)
	}

	log.Log(log.Debug, "no %s match for %s", context, actual)
	return false
}
