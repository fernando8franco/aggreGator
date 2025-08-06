package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	var conf Config
	if err := json.NewDecoder(jsonFile).Decode(&conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name

	jsonConfig, err := json.Marshal(c)
	if err != nil {
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, jsonConfig, os.ModePerm)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := homeDir + "/" + configFileName
	return configFilePath, nil
}
