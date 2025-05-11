package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func Load() (*Config, error) {
	cfg := &Config{}

	configPaths := []string{
		"./config.yaml",
		"./config/config.yaml",
	}

	var foundConfig string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			foundConfig = path
			break
		}
	}

	if foundConfig != "" {
		if err := cleanenv.ReadConfig(foundConfig, cfg); err != nil {
			return nil, err
		}
	} else {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
