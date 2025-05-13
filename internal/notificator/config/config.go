package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type NotificatorPart struct {
	NewIpUrl string `yaml:"new_ip"`
}
type Config struct {
	NotificatorPart `yaml:"notificator"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("cannot load .env file: %w", err)
	}

	yamlConfigPath, ok := os.LookupEnv("YAML_CONFIG_PATH")
	if !ok {
		return nil, errors.New("environment variable YAML_CONFIG_PATH is not set")
	}

	cfg := &Config{}
	if err := cleanenv.ReadConfig(yamlConfigPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot set config: %w", err)
	}

	return cfg, nil
}
