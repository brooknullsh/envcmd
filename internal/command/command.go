package command

import (
  "bufio"
  "fmt"
  "io"
  "os/exec"
  "slices"
  "sync"

  "github.com/brooknullsh/envcmd/internal/config"
  "github.com/brooknullsh/envcmd/internal/kind"
  "github.com/brooknullsh/envcmd/internal/log"
)

func Run(cfg config.Config) {
  var wg sync.WaitGroup

  contents := cfg.Read()
  matchIdx := slices.IndexFunc(contents, func(c config.Config) bool {
    return kind.IsMatch(c.Kind, c.Target)
  })

  if matchIdx == -1 {
    return
  }

  cfg = contents[matchIdx]
  log.Info("\x1b[1m%s\x1b[0m (%s)", cfg.Target, cfg.Kind)

  for idx, cmd := range cfg.Commands {
    if !cfg.Async {
      sharedRun(cmd, idx)
      continue
    }

    wg.Add(1)
    go func() {
      defer wg.Done()
      sharedRun(cmd, idx)
    }()
  }

  wg.Wait()
}

func sharedRun(cmd string, idx int) {
  cmdHandle := exec.Command("bash", "-c", cmd)

  stdout, err := cmdHandle.StdoutPipe()
  if err != nil {
    log.Abort("getting stdout pipe: %v", err)
  }

  stderr, err := cmdHandle.StderrPipe()
  if err != nil {
    log.Abort("getting stderr pipe: %v", err)
  }

  err = cmdHandle.Start()
  if err != nil {
    log.Abort("starting command (%s): %v", cmd, err)
  }

  var wg sync.WaitGroup
  wg.Add(2)

  go printStream(&wg, stdout, idx)
  go printStream(&wg, stderr, idx)

  wg.Wait()

  err = cmdHandle.Wait()
  if err != nil {
    log.Abort("running (%s): %v", cmd, err)
  }
}

func printStream(wg *sync.WaitGroup, stream io.Reader, idx int) {
  defer wg.Done()

  // blue, magenta, cyan, white
  var colours = [4]string{"\033[34m", "\033[35m", "\033[36m", "\033[37m"}
  colourIdx := (idx + 1) % len(colours)
  colour := colours[colourIdx]

  scanner := bufio.NewScanner(stream)
  for scanner.Scan() {
    fmt.Printf("\x1b[1m%s%d\x1b[0m %s\n", colour, idx, scanner.Text())
  }

  if err := scanner.Err(); err != nil {
    log.Warn("reading line: %v", err)
  }
}
