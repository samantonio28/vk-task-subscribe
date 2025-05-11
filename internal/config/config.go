package config

import (
	"time"
)

type Config struct {
	GRPC    GRPCConfig
	Logging LoggingConfig
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env:"GRPC_PORT" env-default:"9000"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"5s"`
}

type LoggingConfig struct {
	FilePath string `yaml:"file_path" env:"LOG_FILE_PATH" env-default:"./logs/access.log"`
	Level    string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}
