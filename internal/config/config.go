package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/brooknullsh/envcmd/internal/log"
)

func configPath() string {
	user, err := user.Current()
	if err != nil {
		log.Log(log.Error, "failed getting user: %v", err)
		os.Exit(1)
	}

	filePath := filepath.Join(user.HomeDir, ".envcmd/config.json")
	return filePath
}

func doesConfigExist() bool {
	filePath := configPath()
	_, err := os.Stat(filePath)
	return err == nil
}

func Create() {
	filePath := configPath()
	dirPath := filepath.Dir(filePath)

	if doesConfigExist() {
		log.Log(log.Info, "exists at %s", filePath)
		return
	}

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Log(log.Error, "creating directory at %s: %v", dirPath, err)
		os.Exit(1)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Log(log.Error, "creating file at %s: %v", filePath, err)
		os.Exit(1)
	}

	defer file.Close()
	log.Log(log.Info, "created %s", filePath)
}

func Delete() {
	filePath := configPath()

	if !doesConfigExist() {
		log.Log(log.Info, "no configuration to delete")
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		log.Log(log.Error, "removing file at %s: %v", filePath, err)
		os.Exit(1)
	}

	log.Log(log.Info, "removed %s", filePath)
}
