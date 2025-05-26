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

func readContent(file *os.File) content {
	var content content
	err := json.NewDecoder(file).Decode(&content)
	if err != nil {
		log.Log(log.Error, "decoding JSON: %v", err)
		os.Exit(1)
	}

	return content
}

func (c *content) WriteToFile(file *os.File) {
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

func (c *Config) doesConfigExist() bool {
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

func (c *Config) Create() {
	if c.doesConfigExist() {
		log.Log(log.Info, "exists at %s", c.filePath)
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

	defaultContent := content{
		Condition: []string{"directory", "foo"},
		Commands:  []string{"echo 'Hello, foo!'"},
	}

	defaultContent.WriteToFile(file)
	defer file.Close()
	log.Log(log.Info, "created %s", c.filePath)
}

func (c *Config) Delete() {
	if !c.doesConfigExist() {
		log.Log(log.Info, "no configuration to delete")
		return
	}

	err := os.Remove(c.filePath)
	if err != nil {
		log.Log(log.Error, "removing file at %s: %v", c.filePath, err)
		os.Exit(1)
	}

	log.Log(log.Info, "removed %s", c.filePath)
}

func (c *Config) List() {
	if !c.doesConfigExist() {
		log.Log(log.Info, "no configuration to list")
		return
	}

	file, err := os.Open(c.filePath)
	if err != nil {
		log.Log(log.Error, "opening file at %s: %v", c.filePath, err)
		os.Exit(1)
	}

	defer file.Close()

	log.Log(log.Debug, "reading %s", c.filePath)
	content := readContent(file)
	log.Log(log.Info, "%v - %v", content.Condition, content.Commands)
}
