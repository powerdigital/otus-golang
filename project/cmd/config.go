package main

import (
	"flag"
	"os"

	"github.com/powerdigital/project/internal/config"
	yaml3 "gopkg.in/yaml.v3"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/tmp/config.yaml", "Path to configuration file")
}

func NewConfig() (*config.Config, error) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config *config.Config
	err = yaml3.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
