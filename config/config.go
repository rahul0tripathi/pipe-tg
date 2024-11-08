package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UID           string `json:"UID" envconfig:"UID"`
	AppID         int    `json:"appID" envconfig:"APP_ID"`
	AppHash       string `json:"appHash" envconfig:"APP_HASH"`
	SessionConfig string `json:"sessionConfig" envconfig:"SESSION_CONFIG"`
	Port          string `json:"port" envconfig:"PORT"`
}

type Option func(*envConfig)

type envConfig struct {
	envFilePath string
}

func defaultEnvConfig() *envConfig {
	return &envConfig{
		envFilePath: ".env",
	}
}

func WithEnvPath(path string) Option {
	return func(cfg *envConfig) {
		cfg.envFilePath = path
	}
}

func LoadEnvConfig(cfg any, options ...Option) error {

	defaultCfg := defaultEnvConfig()
	for _, option := range options {
		option(defaultCfg)
	}

	err := godotenv.Overload(defaultCfg.envFilePath)
	if err != nil {
		return fmt.Errorf("failed to read env file, %w", err)
	}

	err = envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("failed to fill config structure.en, %w", err)
	}

	return nil
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	err := LoadEnvConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
