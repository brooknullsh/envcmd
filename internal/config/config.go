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

func validateCondition(c []string) {
	if len(c) != 2 {
		log.Abort("expected condition length to be 2, is %d", len(c))
	}

	if c[0] != "directory" && c[0] != "branch" {
		log.Abort("condition must be 'directory' or 'branch', is '%s'", c[0])
	}
}

func readContent(f *os.File) []Content {
	var cont []Content

	if err := json.NewDecoder(f).Decode(&cont); err != nil {
		log.Abort("decoding JSON: %v", err)
	}

	for _, c := range cont {
		validateCondition(c.Condition)
	}

	return cont
}

func writeToFile(c []Content, f *os.File) {
	json, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Abort("encoding JSON: %v", err)
	}

	if _, err = f.Write(json); err != nil {
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
	u, err := user.Current()
	if err != nil {
		log.Abort("failed getting user: %v", err)
	}

	c.filePath = filepath.Join(u.HomeDir, ".envcmd/config.json")
}

func (c Config) Create() {
	if c.doesConfigExist() {
		log.Abort("configuration already exists")
	}

	dirPath := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Abort("creating directory at %s: %v", dirPath, err)
	}

	f, err := os.Create(c.filePath)
	if err != nil {
		log.Abort("creating file at %s: %v", c.filePath, err)
	}

	defer f.Close()

	defaultCont := []Content{
		{
			Condition: []string{"directory", "foo"},
			Commands:  []string{"echo 'Hello, foo!'"},
		},
		{
			Condition: []string{"branch", "bar"},
			Commands:  []string{"echo 'Hello, bar!'"},
		},
	}

	writeToFile(defaultCont, f)
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

	f, err := os.Open(c.filePath)
	if err != nil {
		log.Abort("opening file at %s: %v", c.filePath, err)
	}

	defer f.Close()

	log.Log(log.Debug, "read from %s", c.filePath)
	return readContent(f)
}
