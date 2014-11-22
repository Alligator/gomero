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
}

func ReadConfig(path string) Config {
	var config Config
	data, err := ioutil.ReadFile(path)
	err = json.Unmarshal(data, &config)
	err = err
	return config
}
