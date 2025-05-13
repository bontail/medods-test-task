package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type LoggerPart struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output struct {
		Type       string `yaml:"type"`
		Path       string `yaml:"path"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
		Compress   bool   `yaml:"compress"`
	} `yaml:"output"`
}
type Config struct {
	LoggerPart `yaml:"logger"`
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
