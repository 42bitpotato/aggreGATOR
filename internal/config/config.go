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

func Read() (Config, error) {
	jsonFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config json: %v", err)
	}
	return cfg, nil
}

func SetUser(cfg *Config, username string) error {
	cfg.CurrentUserName = username
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = json.NewEncoder(filePath).Encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to set username in config json: %v", err)
	}
	return nil
}
