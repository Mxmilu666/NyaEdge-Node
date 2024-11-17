package source

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type CenterConfig struct {
	Address    string `yaml:"address"`
	Port       int    `yaml:"port"`
	NodeId     string `yaml:"node_id"`
	NodeSecret string `yaml:"node_secret"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Center CenterConfig `yaml:"center"`
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
