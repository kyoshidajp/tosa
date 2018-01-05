package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
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
	path := os.Getenv("HUB_CONFIG")
	if path != "" {
		return path
	}
	path, _ = setHubConfigEnv()
	return path
}

func setHubConfigEnv() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	confPath := filepath.Join(homeDir, ".config", "tosa", "tosa.yml")
	err = os.Setenv("HUB_CONFIG", confPath)
	if err != nil {
		return "", err
	}
	return confPath, nil
}

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
	return host.Browser, nil
}
