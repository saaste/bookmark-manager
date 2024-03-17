package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	confFile   string = "config.yml"
	appVersion string = "0.0.1"
)

type AppConfig struct {
	SiteName          string `yaml:"site_name"`
	Description       string `yaml:"description"`
	BaseURL           string `yaml:"base_url"`
	Password          string `yaml:"password"`
	Secret            string `yaml:"secret"`
	Port              int    `yaml:"port"`
	PageSize          int    `yaml:"page_size"`
	Theme             string `yaml:"theme"`
	CheckInterval     int    `yaml:"check_interval,omitempty"`
	CheckRunOnStartup bool   `yaml:"check_on_app_start,omitempty"`
	GotifyEnabled     bool   `yaml:"gotify_enabled,omitempty"`
	GotifyURL         string `yaml:"gotify_url,omitempty"`
	GotifyToken       string `yaml:"gotify_token,omitempty"`
	AppVersion        string
	AuthorName        string `yaml:"author_name,omitempty"`
	AuthorEmail       string `yaml:"author_email,omitempty"`
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

	appConfig.AppVersion = appVersion

	// Validate base URL
	if !strings.HasSuffix(appConfig.BaseURL, "/") {
		return nil, fmt.Errorf("base URL must have a trailing slash")
	}

	return appConfig, nil
}
