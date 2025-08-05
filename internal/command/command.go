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
  findFirstMatch(&cfg)

  log.Info("\x1b[1m%s\x1b[0m (%s)", cfg.Target, cfg.Kind)
  for idx, cmd := range cfg.Commands {
    if cfg.Async {
      wg.Add(1)
      go sharedRun(&wg, cmd, idx)
    } else {
      sharedRun(&wg, cmd, idx)
    }
  }

  wg.Wait()
}

func findFirstMatch(cfg *config.Config) {
  contents := cfg.Read()
  matchIdx := slices.IndexFunc(contents, func(c config.Config) bool {
    return kind.IsMatch(c.Kind, c.Target)
  })

  if matchIdx == -1 {
    log.Abort("no match found")
  }

  *cfg = contents[matchIdx]
}

func sharedRun(wg *sync.WaitGroup, cmd string, idx int) {
  defer wg.Done()
  cmdHandle := exec.Command("bash", "-c", cmd)

  stdout, err := cmdHandle.StdoutPipe()
  if err != nil {
    log.Abort("getting stdout pipe: %v", err)
  }

  stderr, err := cmdHandle.StderrPipe()
  if err != nil {
    log.Abort("getting stderr pipe: %v", err)
  }

  if err := cmdHandle.Start(); err != nil {
    log.Abort("starting command (%s): %v", cmd, err)
  }

  var wgStream sync.WaitGroup
  wgStream.Add(2)

  go printStream(&wgStream, stdout, idx)
  go printStream(&wgStream, stderr, idx)

  wgStream.Wait()

  if err := cmdHandle.Wait(); err != nil {
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
