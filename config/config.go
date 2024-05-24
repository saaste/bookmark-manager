package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	confFile   string = "config.yml"
	appVersion string = "1.0.2"
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
	CheckerUserAgent  string `yaml:"checker_user_agent,omitempty"`
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

	if appConfig.GotifyEnabled && (appConfig.GotifyToken == "" || appConfig.GotifyURL == "") {
		return nil, fmt.Errorf("gotify notifications are enabled but token or gotify URL is empty")
	}

	// Validate trailing slash of URLS
	if !strings.HasSuffix(appConfig.BaseURL, "/") {
		appConfig.BaseURL = fmt.Sprintf("%s/", appConfig.BaseURL)
	}
	if appConfig.GotifyURL != "" && !strings.HasSuffix(appConfig.GotifyURL, "/") {
		appConfig.GotifyURL = fmt.Sprintf("%s/", appConfig.GotifyURL)
	}

	return appConfig, nil
}
