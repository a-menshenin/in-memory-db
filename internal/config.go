package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Engine struct {
		EngineType string `yaml:"engine_type"`
	} `yaml:"engine"`
	Network struct {
		Address string `yaml:"address"`
		MaxConnections int `yaml:"max_connections",omitempty`
		MaxMessageSize int `yaml:"max_message_size",omitempty`
		IdleTimeout time.Duration `yaml:"idle_timeout",omitempty`
	} `yaml:"network"`
	Logging struct {
		Level string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func GetConfig() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return Config{}, fmt.Errorf("Load env error: %w", err)
   }
	cfgPath := os.Getenv("CONFIG_PATH")

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, fmt.Errorf("Config read error: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Unmarshal config error: %w", err)
	}

	return config, nil
}

func CreateLogger(cfg Config) (*zap.Logger, error) {
	logCfg := zap.NewProductionConfig()
	logCfg.OutputPaths = []string{cfg.Logging.Output}
	logCfg.ErrorOutputPaths = []string{cfg.Logging.Output}
	return logCfg.Build()
}
