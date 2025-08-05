package kind

import (
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "github.com/brooknullsh/envcmd/internal/log"
)

func IsMatch(kind string, target string) bool {
  switch kind {
  case "directory":
    return directoryMatch(target)
  case "branch":
    return branchMatch(target)
  default:
    log.Warn("unknown kind (%s)", kind)
    return false
  }
}

func directoryMatch(target string) bool {
  dirPath, err := os.Getwd()
  if err != nil {
    log.Abort("reading working directory: %v", err)
  }

  return filepath.Base(dirPath) == target
}

func branchMatch(target string) bool {
  cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

  out, err := cmd.Output()
  if err != nil {
    if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 128 {
      log.Warn("no git in working directory")
      return false
    }

    log.Abort("reading branch: %v", err)
  }

  return strings.TrimSpace(string(out)) == target
}
