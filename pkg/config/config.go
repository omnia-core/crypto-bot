package config

import (
	"os"

	"crypto-bot/pkg/log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env   string      `yaml:"env"`
	Upbit UpbitConfig `yaml:"upbit"`
}

type UpbitConfig struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
}

func ParseConfig() *Config {
	var cfg Config
	configFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.New().Fatalf("err read local config yaml %v", err)
	}
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		log.New().Fatalf("err unmarshal local config yaml %v", err)
	}
	return &cfg
}
