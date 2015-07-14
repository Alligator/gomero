package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Name     string
	Host     string
	Channels []string
	Prefix   string
	ApiKeys  map[string]string
}

func ReadConfig(path string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(path)
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
