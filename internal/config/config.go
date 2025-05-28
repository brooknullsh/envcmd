package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/brooknullsh/envcmd/internal/log"
)

type Content struct {
	Condition []string `json:"condition"`
	Commands  []string `json:"commands"`
}

func readContent(file *os.File) []Content {
	var content []Content

	err := json.NewDecoder(file).Decode(&content)
	if err != nil {
		log.Log(log.Error, "decoding JSON: %v", err)
		os.Exit(1)
	}

	return content
}

func writeToFile(c []Content, file *os.File) {
	jsonContent, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Log(log.Error, "encoding JSON: %v", err)
		os.Exit(1)
	}

	_, err = file.Write(jsonContent)
	if err != nil {
		log.Log(log.Error, "writing: %v", err)
		os.Exit(1)
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
		log.Log(log.Error, "failed getting user: %v", err)
		os.Exit(1)
	}

	c.filePath = filepath.Join(user.HomeDir, ".envcmd/config.json")
}

func (c Config) Create() {
	if c.doesConfigExist() {
		log.Log(log.Error, "configuration already exists")
		return
	}

	dirPath := filepath.Dir(c.filePath)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Log(log.Error, "creating directory at %s: %v", dirPath, err)
		os.Exit(1)
	}

	file, err := os.Create(c.filePath)
	if err != nil {
		log.Log(log.Error, "creating file at %s: %v", c.filePath, err)
		os.Exit(1)
	}

	defer file.Close()

	defaultContent := []Content{
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
		log.Log(log.Error, "configuration doesn't exist")
		return
	}

	err := os.Remove(c.filePath)
	if err != nil {
		log.Log(log.Error, "removing file at %s: %v", c.filePath, err)
		os.Exit(1)
	}

	log.Log(log.Info, "deleted from %s", c.filePath)
}

func (c Config) Read() []Content {
	if !c.doesConfigExist() {
		log.Log(log.Error, "configuration doesn't exist")
		return nil
	}

	file, err := os.Open(c.filePath)
	if err != nil {
		log.Log(log.Error, "opening file at %s: %v", c.filePath, err)
		os.Exit(1)
	}

	defer file.Close()

	log.Log(log.Debug, "read from %s", c.filePath)
	return readContent(file)
}
