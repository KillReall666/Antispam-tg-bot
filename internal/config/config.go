package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const configFile = ".o/config.yaml"

type Config struct {
	Token string `yaml:"token"` // Тг апи токен.
}

type Service struct {
	config Config
}

func New() (*Service, error) {
	s := &Service{}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("error reading config file:", err)
		return nil, errors.New(fmt.Sprintf("error reading config file: %s", err))
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		log.Println("error parsing yaml file:", err)
		return nil, errors.New(fmt.Sprintf("error parsing yaml file: %s", err))
	}
	return s, nil
}

func (s *Service) Token() string {
	return s.config.Token
}

func (s *Service) GetConfig() Config {
	return s.config
}
