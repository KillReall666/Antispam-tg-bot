package appcfg

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const configFile = "./appconfig.yaml"

type AppConfig struct {
	Token     string `yaml:"token"` // Тг апи токен.
	DBConnStr string `yaml:"db_conn_str"`
	TgApiURL  string `yaml:"tg_api_url"`
}

func New() (*AppConfig, error) {
	cfg := &AppConfig{}

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

func (c *AppConfig) GetToken() string {
	return c.Token
}

func (c *AppConfig) GetConfig() AppConfig {
	return AppConfig{}
}
