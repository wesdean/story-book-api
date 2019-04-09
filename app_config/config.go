package app_config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Models conifigModels
}

type conifigModels struct {
	DatabaseSeed string
}

func NewConfigFromFile(filename string) (*Config, error) {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (config *Config) Validate() error {
	return nil
}
