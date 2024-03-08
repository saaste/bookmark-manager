package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	confFile string = "config.yml"
)

type AppConfig struct {
	SiteName string `yaml:"site_name"`
	BaseURL  string `yaml:"base_url"`
	Password string `yaml:"password"`
	Secret   string `yaml:"secret"`
	Port     int    `yaml:"port"`
	PageSize int    `yaml:"page_size"`
	Template string `yaml:"template"`
}

func LoadConfig() (*AppConfig, error) {
	yamlFile, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{}
	err = yaml.Unmarshal(yamlFile, appConfig)
	if err != nil {
		return nil, err
	}

	return appConfig, nil
}
