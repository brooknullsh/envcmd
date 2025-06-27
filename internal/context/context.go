package context

import (
	"os"
	"os/exec"
	"path/filepath"
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
	switch content.Context {
	case "directory":
		return directoryMatch(content.Target)
	case "branch":
		return branchMatch(content.Target)
	default:
		log.Warn("unknown context '%s' in %s", content.Context, content.Name)
		return false
	}
}
