package context

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/brooknullsh/envcmd/internal/config"
	"github.com/brooknullsh/envcmd/internal/log"
)

func branchMatch(target string) bool {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	out, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 128 {
			log.Warn("no git in current directory")
			return false
		}
		log.Abort("reading git branch -> %v", err)
	}

	return strings.TrimSpace(string(out)) == target
}

func directoryMatch(target string) bool {
	dirPath, err := os.Getwd()
	if err != nil {
		log.Abort("reading current directory -> %v", err)
	}

	return filepath.Base(dirPath) == target
}

func Match(content config.Content) bool {
	var matched bool

	switch content.Context {
	case "directory":
		matched = slices.ContainsFunc(content.Targets, func(target string) bool {
			return directoryMatch(target)
		})
	case "branch":
		matched = slices.ContainsFunc(content.Targets, func(target string) bool {
			return branchMatch(target)
		})
	case "both":
		if len(content.Targets) != 2 {
			log.Abort("the 'both' context should be 2 in length -> %s", content.Name)
		}

		matched = directoryMatch(content.Targets[0]) && branchMatch(content.Targets[1])
	default:
		log.Warn("unknown context '%s' in %s", content.Context, content.Name)
		return false
	}

	return matched
}
