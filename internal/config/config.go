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

  path string
}

func (c *Config) Init() {
  user, err := user.Current()
  if err != nil {
    log.Abort("finding user: %v", err)
  }

  c.path = filepath.Join(user.HomeDir, ".envcmd/config.json")
}

func (c Config) Create() {
  if c.exists() {
    log.Abort("already exists: %s", c.path)
  }

  dirPath := filepath.Dir(c.path)
  if err := os.Mkdir(dirPath, 0755); err != nil {
    log.Abort("creating directory (%s): %v", dirPath, err)
  }

  file, err := os.Create(c.path)
  if err != nil {
    log.Abort("creating file (%s): %v", c.path, err)
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

  encodeAndWriteJSON(contents, c.path)
  log.Info("created %s", c.path)
}

func (c Config) Delete() {
  if !c.exists() {
    log.Abort("no configuration found")
  }

  log.Warn("delete (y/N): %s", c.path)
  answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
  if err != nil || strings.TrimSpace(answer) != "y" {
    return
  }

  if err := os.Remove(c.path); err != nil {
    log.Abort("removing file (%s): %v", c.path, err)
  }

  dirPath := filepath.Dir(c.path)
  if err := os.Remove(dirPath); err != nil {
    log.Abort("removing directory (%s): %v", dirPath, err)
  }

  log.Info("deleted %s", c.path)
}

func (c Config) List() {
  if !c.exists() {
    log.Abort("no configuration found")
  }

  for idx, cfg := range c.Read() {
    log.Info(
      "%d \x1b[1m%s\033[0m (%s) (%s)",
      idx,
      cfg.Target,
      cfg.Kind,
      cfg.asyncToString(),
    )

    for _, cmd := range cfg.Commands {
      log.Info("$ %s", cmd)
    }
  }
}

func (c Config) Read() []Config {
  if !c.exists() {
    log.Abort("no configuration found")
  }

  return configFromFile(c.path)
}

func (c Config) Add(newCfg Config) {
  if !c.exists() {
    log.Abort("no configuration found")
  }

  contents := configFromFile(c.path)
  contents = append(contents, newCfg)

  encodeAndWriteJSON(contents, c.path)
  log.Info(
    "added: \x1b[1m%s\033[0m (%s) (%v)",
    newCfg.Target,
    newCfg.Kind,
    newCfg.asyncToString(),
  )
}

func (c Config) exists() bool {
  _, err := os.Stat(c.path)
  return err == nil
}

func (c Config) asyncToString() string {
  if c.Async {
    return "async"
  } else {
    return "sync"
  }
}

func encodeAndWriteJSON(contents []Config, path string) {
  json, err := json.MarshalIndent(contents, "", "  ")
  if err != nil {
    log.Abort("encoding JSON: %v", err)
  }

  if err := os.WriteFile(path, json, 0644); err != nil {
    log.Abort("writing JSON: %v", err)
  }
}

func configFromFile(path string) []Config {
  file, err := os.Open(path)
  if err != nil {
    log.Abort("opening file (%s): %v", path, err)
  }

  defer file.Close()

  var contents []Config
  if err := json.NewDecoder(file).Decode(&contents); err != nil {
    log.Abort("decoding JSON: %v", err)
  }

  return contents
}
