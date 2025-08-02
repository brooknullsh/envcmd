package config

import (
  "bufio"
  "encoding/json"
  "os"
  "os/user"
  "path/filepath"
  "strings"

  "github.com/brooknullsh/envcmd/internal/log"
)

type Config struct {
  Async    bool     `json:"async"`
  Kind     string   `json:"kind"`
  Target   string   `json:"target"`
  Commands []string `json:"commands"`

  filePath string
}

func (c *Config) Init() {
  user, err := user.Current()
  if err != nil {
    log.Abort("finding user: %v", err)
  }

  c.filePath = filepath.Join(user.HomeDir, ".envcmd/config.json")
}

func (c Config) Create() {
  if c.exists() {
    log.Abort("already exists: %s", c.filePath)
  }

  dirPath := filepath.Dir(c.filePath)
  if err := os.Mkdir(dirPath, 0755); err != nil {
    log.Abort("creating directory (%s): %v", dirPath, err)
  }

  file, err := os.Create(c.filePath)
  if err != nil {
    log.Abort("creating file (%s): %v", c.filePath, err)
  }

  defer file.Close()
  contents := []Config{
    {
      Async:    false,
      Kind:     "directory",
      Target:   "foo",
      Commands: []string{"echo 'Hello, world!'"},
    },
  }

  json, err := json.MarshalIndent(contents, "", "  ")
  if err != nil {
    log.Abort("encoding JSON: %v", err)
  }

  if _, err := file.Write(json); err != nil {
    log.Abort("writing JSON: %v", err)
  }

  log.Info("created %s", c.filePath)
}

func (c Config) Delete() {
  if !c.exists() {
    log.Abort("no configuration found")
  }

  log.Warn("delete (y/N): %s", c.filePath)
  answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
  if err != nil || strings.TrimSpace(answer) != "y" {
    return
  }

  if err := os.Remove(c.filePath); err != nil {
    log.Abort("removing file (%s): %v", c.filePath, err)
  }

  dirPath := filepath.Dir(c.filePath)
  if err := os.Remove(dirPath); err != nil {
    log.Abort("removing directory (%s): %v", dirPath, err)
  }

  log.Info("deleted %s", c.filePath)
}

func (c Config) List() {
  if !c.exists() {
    log.Abort("no configuration found")
  }

  for idx, cfg := range c.Read() {
    var isAsync string
    if cfg.Async {
      isAsync = "async"
    } else {
      isAsync = "sync"
    }

    log.Info("%d \x1b[1m%s\033[0m (%s) (%s)", idx, cfg.Target, cfg.Kind, isAsync)
    for _, cmd := range cfg.Commands {
      log.Info("$ %s", cmd)
    }
  }
}

func (c Config) Read() []Config {
  if !c.exists() {
    log.Abort("no configuration found")
  }

  file, err := os.Open(c.filePath)
  if err != nil {
    log.Abort("opening file (%s): %v", c.filePath, err)
  }

  defer file.Close()

  var contents []Config
  if err := json.NewDecoder(file).Decode(&contents); err != nil {
    log.Abort("decoding JSON: %v", err)
  }

  return contents
}

func (c Config) exists() bool {
  _, err := os.Stat(c.filePath)
  return err == nil
}
