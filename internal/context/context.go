package context

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/brooknullsh/envcmd/internal/log"
)

func branchMatch(actual string) bool {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	output, err := cmd.Output()
	if err != nil {
		log.Abort("reading git branch: %v", err)
	}

	return strings.TrimSpace(string(output)) == actual
}

func directoryMatch(actual string) bool {
	path, err := os.Getwd()
	if err != nil {
		log.Abort("reading current directory: %v", err)
	}

	return filepath.Base(path) == actual
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
