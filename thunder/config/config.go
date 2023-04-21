package config

import (
	"math/big"
	"os"

	"github.com/ledgerwatch/erigon-lib/thunder/utils"
	"gopkg.in/yaml.v2"
)

type Key struct {
	GenesisCommPath string `toml:"GenesisCommPath" yaml:"GenesisCommPath"`
	AlterCommPath   string `toml:"alterCommPath" yaml:"alterCommPath"`
}

type Pala struct {
	FromGenesis bool `toml:"fromGenesis" yaml:"fromGenesis"`
}

type Chain struct {
	InitialSupply      string   `toml:"initialSupply" yaml:"initialSupply"`
	InitialSupplyInBig *big.Int `toml:"-" yaml:"-"`
}

type Config struct {
	Key   Key   `toml:"key" yaml:"key"`
	Pala  Pala  `toml:"pala" yaml:"pala"`
	Chain Chain `toml:"chain" yaml:"chain"`
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

	cfg.Chain.InitialSupplyInBig, err = utils.SimpleScientificBigIntParse(cfg.Chain.InitialSupply)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
