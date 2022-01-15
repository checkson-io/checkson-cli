package config

import (
	"os"
	"path/filepath"
)

func GetConfigPath() string {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configPath := filepath.Join(homeDir, ".checkson")
	return configPath
}
