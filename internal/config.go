package internal

import (
	"fmt"
	"os"
	"time"

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
		ServerLevel string `yaml:"server_level"`
		CliOutput string `yaml:"cli_output"`
		ServerOutput string `yaml:"server_output"`
	} `yaml:"logging"`
}

func GetConfig() Config {
	// cfgPath, exists := os.LookupEnv("CONFIG_PATH")
	// if !exists {
	// 	fmt.Println("Config read error")

	// 	os.Exit(1)
	// }

	data, err := os.ReadFile("/Users/menshenin/GolandProjects/BalunProjects/InMemoryDB2/config.yaml")
	if err != nil {
		fmt.Println("Config read error")
		
		return Config{}
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Unmarshal config error")

		return Config{}
	}

	return config
}

func CreateLogger(cfg Config) (*zap.Logger, error) {
	logCfg := zap.NewProductionConfig()
	logCfg.OutputPaths = []string{cfg.Logging.CliOutput}
	logCfg.ErrorOutputPaths = []string{cfg.Logging.CliOutput}
	return logCfg.Build()
}
