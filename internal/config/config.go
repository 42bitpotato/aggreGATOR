package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/42bitpotato/aggreGATOR/internal/database"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
	DateFormat      string `json:"date_format"`
}

type State struct {
	Db     *database.Queries
	Cfg    *Config
	Logger *log.Logger
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

	// Set date format
	if cfg.DateFormat == "" {
		cfg.DateFormat = "2006-01-02 15:04:05"
		err = write(&cfg)
		if err != nil {
			return Config{}, fmt.Errorf("error writing date to config json: %v", err)
		}
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
