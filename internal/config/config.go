package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/brooknullsh/envcmd/internal/log"
)

type Content struct {
	Name     string   `json:"name"`
	Async    bool     `json:"async"`
	Context  string   `json:"context"`
	Targets  []string `json:"targets"`
	Commands []string `json:"commands"`
}

func (c Content) Print() {
	log.Info("Name     \x1b[1m%s\033[0m", c.Name)
	log.Info("Context  \x1b[1m%s\033[0m", c.Context)
	log.Info("Targets  \x1b[1m%v\033[0m", c.Targets)
	log.Info("Async    \x1b[1m%v\033[0m", c.Async)

	fmt.Println()
	for _, cmd := range c.Commands {
		log.Info("\x1b[1m%s\033[0m", cmd)
	}
}

func readContent(file *os.File) []Content {
	var contents []Content

	if err := json.NewDecoder(file).Decode(&contents); err != nil {
		log.Abort("decoding JSON -> %v", err)
	}

	return contents
}

func writeToFile(contents []Content, file *os.File) {
	json, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		log.Abort("encoding JSON -> %v", err)
	}

	if _, err = file.Write(json); err != nil {
		log.Abort("writing -> %v", err)
	}
}

type Config struct {
	filePath string
}

func (c Config) configExists() bool {
	_, err := os.Stat(c.filePath)
	return err == nil
}

func (c *Config) InitPath() {
	user, err := user.Current()
	if err != nil {
		log.Abort("getting user -> %v", err)
	}

	c.filePath = filepath.Join(user.HomeDir, ".envcmd/config.json")
}

func (c Config) Create() {
	if c.configExists() {
		log.Abort("already exists -> %s", c.filePath)
	}

	dirPath := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Abort("creating directory -> %s: %v", dirPath, err)
	}

	file, err := os.Create(c.filePath)
	if err != nil {
		log.Abort("creating file -> %s: %v", c.filePath, err)
	}

	defer file.Close()
	content := []Content{
		{
			Name:     "foo",
			Async:    true,
			Context:  "directory",
			Targets:  []string{"bar"},
			Commands: []string{"echo 'Hello, bar!'"},
		},
	}

	writeToFile(content, file)
	log.Info("created -> %s", c.filePath)
}

func (c Config) Delete() {
	if !c.configExists() {
		log.Abort("no config found")
	}

	if err := os.Remove(c.filePath); err != nil {
		log.Abort("removing file -> %s: %v", c.filePath, err)
	}

	dirPath := filepath.Dir(c.filePath)
	if err := os.Remove(dirPath); err != nil {
		log.Abort("removing directory -> %s: %v", dirPath, err)
	}

	log.Info("deleted -> %s", c.filePath)
}

func (c Config) Read() []Content {
	if !c.configExists() {
		log.Abort("no config found")
	}

	file, err := os.Open(c.filePath)
	if err != nil {
		log.Abort("opening file -> %s: %v", c.filePath, err)
	}

	defer file.Close()
	return readContent(file)
}
