package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/brooknullsh/envcmd/internal/log"
)

type Content struct {
	Async     bool     `json:"async"`
	Condition []string `json:"condition"`
	Commands  []string `json:"commands"`
}

func validateCondition(condition []string) {
	if len(condition) != 2 {
		log.Abort("expected condition length to be 2, is %d", len(condition))
	}

	if condition[0] != "directory" && condition[0] != "branch" {
		log.Abort("condition must be 'directory' or 'branch', is '%s'", condition[0])
	}
}

func readContent(configFile *os.File) []Content {
	var configContent []Content

	if err := json.NewDecoder(configFile).Decode(&configContent); err != nil {
		log.Abort("decoding JSON: %v", err)
	}

	for _, content := range configContent {
		validateCondition(content.Condition)
	}

	return configContent
}

func writeToFile(contents []Content, configFile *os.File) {
	json, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		log.Abort("encoding JSON: %v", err)
	}

	if _, err = configFile.Write(json); err != nil {
		log.Abort("writing: %v", err)
	}
}

type Config struct {
	filePath string
}

func (c Config) doesConfigExist() bool {
	_, err := os.Stat(c.filePath)
	return err == nil
}

func (c *Config) InitPath() {
	user, err := user.Current()
	if err != nil {
		log.Abort("failed getting user: %v", err)
	}

	c.filePath = filepath.Join(user.HomeDir, ".envcmd/config.json")
}

func (c Config) Create() {
	if c.doesConfigExist() {
		log.Abort("configuration already exists")
	}

	dirPath := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Abort("creating directory at %s: %v", dirPath, err)
	}

	configFile, err := os.Create(c.filePath)
	if err != nil {
		log.Abort("creating file at %s: %v", c.filePath, err)
	}

	defer configFile.Close()

	defaultContent := []Content{
		{
			Async:     true,
			Condition: []string{"directory", "foo"},
			Commands:  []string{"echo 'Hello, foo!'"},
		},
		{
			Async:     false,
			Condition: []string{"branch", "bar"},
			Commands:  []string{"echo 'Hello, bar!'"},
		},
	}

	writeToFile(defaultContent, configFile)
	log.Log(log.Info, "created at %s", c.filePath)
}

func (c Config) Delete() {
	if !c.doesConfigExist() {
		log.Abort("configuration doesn't exist")
	}

	if err := os.Remove(c.filePath); err != nil {
		log.Abort("removing file at %s: %v", c.filePath, err)
	}

	dirPath := filepath.Dir(c.filePath)
	if err := os.Remove(dirPath); err != nil {
		log.Abort("removing directory at %s: %v", dirPath, err)
	}

	log.Log(log.Info, "deleted from %s", c.filePath)
}

func (c Config) Read() []Content {
	if !c.doesConfigExist() {
		log.Abort("configuration doesn't exist")
	}

	configFile, err := os.Open(c.filePath)
	if err != nil {
		log.Abort("opening file at %s: %v", c.filePath, err)
	}

	defer configFile.Close()
	return readContent(configFile)
}
