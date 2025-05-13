package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type ServerPart struct {
	GinMode                    string `env:"GIN_MODE" envDefault:"debug"`
	Host                       string `yaml:"host" env-default:"localhost"`
	Port                       int    `yaml:"port" env-default:"5432"`
	AccessTokenLifetimeSeconds int    `yaml:"access_token_lifetime_seconds"`
	RefreshTokenLifetimeHours  int    `yaml:"refresh_token_lifetime_hours"`
	DatabaseUrl                string `env:"DATABASE_URL" env-required:"true"`
	SecretKey                  string `env:"SECRET_KEY" env-required:"true"`
}
type Config struct {
	ServerPart `yaml:"server"`
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
