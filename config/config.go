package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"database"`

	Jwt struct {
		Secret     string `yaml:"secret"`
		Iss        string `yaml:"iss"`
		TokenAge   int    `yaml:"tokenAge"`
		RefreshAge int    `yaml:"refreshAge"`
	} `yaml:"jwt"`
}

func GetConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
