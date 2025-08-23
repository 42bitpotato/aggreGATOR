package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	Cfg *Config
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home dir: %v", err)
	}
	filePath := filepath.Join(homeDir, configFileName)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	} else if os.IsNotExist(err) {
		return "", fmt.Errorf("unable to get config json: %v", err)
	}
	return "", fmt.Errorf("error getting config file path: %v", err)
}

func write(cfg *Config) error {
	jsonFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	err = json.NewEncoder(jsonFile).Encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to write config to json: %v", err)
	}
	return nil
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
	defer jsonFile.Close()

	var cfg Config
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("unable to decode config json: %v", err)
	}
	return cfg, nil
}

func SetUser(cfg *Config, username string) error {
	cfg.CurrentUserName = username
	err := write(cfg)
	if err != nil {
		return fmt.Errorf("failed to write username to config json: %v", err)
	}
	return nil
}
