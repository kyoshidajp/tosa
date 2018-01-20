package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type yamlHost struct {
	User       string `yaml:"user"`
	OAuthToken string `yaml:"oauth_token"`
	Protocol   string `yaml:"protocol"`
	Browser    string `yaml:"browser"`
}

type yamlConfig map[string][]yamlHost

func configsFile() string {
	return os.Getenv("HUB_CONFIG")
}

// GetBrowser gets name of browser
func GetBrowser() (string, error) {
	buf, err := ioutil.ReadFile(configsFile())
	if err != nil {
		return "", err
	}

	yc := make(yamlConfig)
	err = yaml.Unmarshal(buf, &yc)
	if err != nil {
		return "", err
	}

	host := yc["github.com"][0]
	Debugf("Config: %v", host)

	return host.Browser, nil
}
