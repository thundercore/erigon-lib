package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Key struct {
	GenesisCommPath string `toml:"GenesisCommPath" yaml:"GenesisCommPath"`
}

type Config struct {
	Key Key `toml:"key" yaml:"key"`
}

func New(path string) (*Config, error) {
	cfg := &Config{}

	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buffer, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
