package context

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/brooknullsh/envcmd/internal/log"
)

func branchMatch(exp string) bool {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	out, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok && e.ExitCode() == 128 {
			log.Log(log.Debug, "no git in current directory")
			return false
		}

		log.Abort("reading git branch: %v", err)
	}

	return strings.TrimSpace(string(out)) == exp
}

func directoryMatch(exp string) bool {
	p, err := os.Getwd()
	if err != nil {
		log.Abort("reading current directory: %v", err)
	}

	return filepath.Base(p) == exp
}

func Match(ctx string, exp string) bool {
	if ctx == "directory" {
		return directoryMatch(exp)
	} else if ctx == "branch" {
		return branchMatch(exp)
	}

	log.Log(log.Debug, "no %s match for %s", ctx, exp)
	return false
}
