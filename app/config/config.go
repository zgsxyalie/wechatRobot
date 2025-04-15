package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type App struct {
	SdkPath string `yaml:"sdkPath"`
}

type apiConfig struct {
	URL         string  `yaml:"url"`
	Key         string  `yaml:"key"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
}

type aiConfig struct {
	Prompt  string `yaml:"prompt"`
	Trigger string `yaml:"trigger"`
}

type config struct {
	App   App       `yaml:"app"`
	Debug bool      `yaml:"debug"`
	API   apiConfig `yaml:"api"`
	AI    aiConfig  `yaml:"ai"`
}

var AppConfig *config

func init() {
	AppConfig = loadConfig()
}

func loadConfig() *config {
	configFile := "config.yaml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file %s: %v", configFile, err)
	}

	var cfg config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &cfg
}
