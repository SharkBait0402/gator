package config

import (
	"os"
	"encoding/json"
	"path/filepath"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	path := filepath.Join(home, configFileName)

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err!=nil {
		return Config{}, err
	}

	return cfg, nil

}

func (cfg *Config) SetUser(name string) error {

	cfg.CurrentUserName = name

	home, err := os.UserHomeDir()
	if err !=nil {
		return err
	}

	path := filepath.Join(home,  configFileName)

	data, err := json.Marshal(*cfg)
	if err!=nil {
		return err
	}

	err = os.WriteFile(path, data, 0o644)
	if err!=nil {
		return err
	}

	return nil

}
