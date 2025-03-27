package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DBUrl 					string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	cfg := Config{}

	path, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name

	byteValue, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, byteValue, 0644)
}