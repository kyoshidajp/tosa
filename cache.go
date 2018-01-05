package main

import (
	"io/ioutil"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Sha struct {
	HTMLUrl string
}

type Cache struct {
	ShaMap map[string]*Sha
}

type yamlCache struct {
	Sha     string `yaml:"sha"`
	HTMLUrl string `yaml:"html_url"`
}

func cacheFile() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeDir, ".config", "tosa", "cache.yml")
	if err != nil {
		return "", err
	}

	return path, nil
}

func NewCache() (*Cache, error) {
	file, err := cacheFile()
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	yc := make(map[string][]yamlCache)
	err = yaml.Unmarshal(buf, &yc)
	if err != nil {
		return nil, err
	}

	caches := yc["github.com"]
	shaMap := make(map[string]*Sha)
	for _, c := range caches {
		shaMap[c.Sha] = &Sha{HTMLUrl: c.HTMLUrl}
	}

	return &Cache{
		ShaMap: shaMap,
	}, nil
}

func (c *Cache) GetHTMLUrl(sha string) (string, error) {
	a, exist := c.ShaMap[sha]
	if !exist {
		return "", nil
	}
	return a.HTMLUrl, nil
}
