package gigachatcfg

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type GigachatConfig struct {
	ClientID     string `yaml:"clientID"`
	Scope        string `yaml:"scope"`
	ClientSecret string `yaml:"clientSecret"`
	AuthDataHash string `yaml:"authDataHash"`
	AuthUrl      string `yaml:"authUrl"`
	RequestUrl   string `yaml:"requestUrl"`
	RedisAddr    string `yaml:"redis_addr"`
}

const configFile = "./gigachatconfig.yaml"

func New() (*GigachatConfig, error) {
	cfg := &GigachatConfig{}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("error reading config file:", err)
		return nil, errors.New(fmt.Sprintf("error reading config file: %s", err))
	}

	err = yaml.Unmarshal(rawYAML, &cfg)
	if err != nil {
		log.Println("error parsing yaml file:", err)
		return nil, errors.New(fmt.Sprintf("error parsing yaml file: %s", err))
	}
	return cfg, nil
}
