package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/brooknullsh/envcmd/internal/log"
)

type Schema struct {
	Name     string   `json:"name"`
	Async    bool     `json:"async"`
	Context  string   `json:"context"`
	Targets  []string `json:"targets"`
	Commands []string `json:"commands"`
}

type Config struct {
	filePath string
}

func (c Config) configExists() bool {
	_, err := os.Stat(c.filePath)
	return err == nil
}

func (c *Config) Init() {
	user, err := user.Current()
	if err != nil {
		log.Abort("getting user -> %v", err)
	}

	c.filePath = filepath.Join(user.HomeDir, ".envcmd/config.json")
}

func (c Config) Read() []Schema {
	if !c.configExists() {
		log.Abort("no config found")
	}

	file, err := os.Open(c.filePath)
	if err != nil {
		log.Abort("opening file -> %s: %v", c.filePath, err)
	}

	defer file.Close()

	var contents []Schema
	if err := json.NewDecoder(file).Decode(&contents); err != nil {
		log.Abort("decoding JSON -> %v", err)
	}

	return contents
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
	contents := []Schema{
		{
			Name:     "foo",
			Async:    true,
			Context:  "directory",
			Targets:  []string{"bar"},
			Commands: []string{"echo 'Hello, bar!'"},
		},
	}

	json, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		log.Abort("encoding JSON -> %v", err)
	}

	if _, err = file.Write(json); err != nil {
		log.Abort("writing -> %v", err)
	}

	log.Info("created -> %s", c.filePath)
}

func (c Config) Delete() {
	if !c.configExists() {
		log.Abort("no config found")
	}

	log.Warn("delete -> %s (y/N)", c.filePath)
	res, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	if strings.TrimSpace(res) != "y" {
		return
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

func (c Config) List() {
	if !c.configExists() {
		log.Abort("no config found")
	}

	for _, obj := range c.Read() {
		fmt.Println()

		log.Info("Name     \x1b[1m%s\033[0m", obj.Name)
		log.Info("Context  \x1b[1m%s\033[0m", obj.Context)
		log.Info("Targets  \x1b[1m%v\033[0m", obj.Targets)
		log.Info("Async    \x1b[1m%v\033[0m", obj.Async)

		fmt.Println()
		for _, cmd := range obj.Commands {
			fmt.Printf(">> \x1b[1m%s\033[0m\n", cmd)
		}
	}
}
