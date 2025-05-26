package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/brooknullsh/envcmd/internal/cmd"
	"github.com/brooknullsh/envcmd/internal/log"
)

type content struct {
	Condition []string `json:"condition"`
	Commands  []string `json:"commands"`
}

func (c *content) Process(run bool) {
	for _, command := range c.Commands {
		log.Log(
			log.Debug,
			"if \x1b[1m%v\033[0m is \x1b[1m%v\033[0m then run \x1b[1m%v",
			c.Condition[0],
			c.Condition[1],
			command,
		)

		if run {
			cmd.RunCmd(command)
		}
	}
}

func readContent(file *os.File) []content {
	var content []content

	err := json.NewDecoder(file).Decode(&content)
	if err != nil {
		log.Log(log.Error, "decoding JSON: %v", err)
		os.Exit(1)
	}

	return content
}

func writeToFile(c []content, file *os.File) {
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
}

func (c *Config) Delete() {
	if !c.doesConfigExist() {
		return
	}

	err := os.Remove(c.filePath)
	if err != nil {
		log.Log(log.Error, "removing file at %s: %v", c.filePath, err)
		os.Exit(1)
	}
}

func (c *Config) Read() []content {
	if !c.doesConfigExist() {
		return nil
	}

	file, err := os.Open(c.filePath)
	if err != nil {
		log.Log(log.Error, "opening file at %s: %v", c.filePath, err)
		os.Exit(1)
	}

	defer file.Close()
	return readContent(file)
}
