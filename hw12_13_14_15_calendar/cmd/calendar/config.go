package main

import (
	"io/ioutil"

	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/config"
	"gopkg.in/yaml.v3"
)

func NewConfig() (*config.Config, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config *config.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
