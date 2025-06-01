package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/brooknullsh/envcmd/internal/log"
)

type content struct {
	Condition []string `json:"condition"`
	Commands  []string `json:"commands"`
}

func validateCondition(c []string) {
	if len(c) != 2 {
		log.Abort("expected condition length to be 2, is %d", len(c))
	}

	if c[0] != "directory" && c[0] != "branch" {
		log.Abort("condition must be 'directory' or 'branch', is '%s'", c[0])
	}
}

func readContent(file *os.File) []content {
	var content []content

	if err := json.NewDecoder(file).Decode(&content); err != nil {
		log.Abort("decoding JSON: %v", err)
	}

	for _, c := range content {
		validateCondition(c.Condition)
	}

	return content
}

func writeToFile(c []content, file *os.File) {
	jsonContent, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Abort("encoding JSON: %v", err)
	}

	if _, err = file.Write(jsonContent); err != nil {
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

func (c *Config) FindPath() {
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

	file, err := os.Create(c.filePath)
	if err != nil {
		log.Abort("creating file at %s: %v", c.filePath, err)
	}

	defer file.Close()

	defaultContent := []content{
		{
			Condition: []string{"directory", "foo"},
			Commands:  []string{"echo 'Hello, foo!'"},
		},
		{
			Condition: []string{"branch", "bar"},
			Commands:  []string{"echo 'Hello, bar!'"},
		},
	}

	writeToFile(defaultContent, file)
	log.Log(log.Info, "created at %s", c.filePath)
}

func (c Config) Delete() {
	if !c.doesConfigExist() {
		log.Abort("configuration doesn't exist")
	}

	if err := os.Remove(c.filePath); err != nil {
		log.Abort("removing file at %s: %v", c.filePath, err)
	}

	log.Log(log.Info, "deleted from %s", c.filePath)
}

func (c Config) Read() []content {
	if !c.doesConfigExist() {
		log.Abort("configuration doesn't exist")
	}

	file, err := os.Open(c.filePath)
	if err != nil {
		log.Abort("opening file at %s: %v", c.filePath, err)
	}

	defer file.Close()

	log.Log(log.Debug, "read from %s", c.filePath)
	return readContent(file)
}
