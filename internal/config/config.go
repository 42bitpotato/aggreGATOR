package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home dir: %v", err)
	}
	filePath := filepath.Join(homeDir, configFileName)
	if _, err := os.Stat(filePath); err == nil {
		return homeDir, nil
	} else if os.IsNotExist(err) {
		return "", fmt.Errorf("unable to get config json: %v", err)
	}
	return "", fmt.Errorf("error getting config file path: %v", err)
}

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func ReadConfigJson() (Config, error) {
	jsonFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	var config Config
	decoder = json.NewDecoder(jsonFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config json: %v", err)
	}
	return config, nil
}
