package config

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host string `envconfig:"HOST"`
	Port string `envconfig:"PORT"`
}

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		ReRead()
	})

	return &config
}

func ReRead() *Config {
	log.Println("reading app config")
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	_, err = json.MarshalIndent(config, "", "")
	configBytes, err := json.MarshalIndent(config, "", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Configuration", string(configBytes))
	return &config
}

func WriteConfig(stringWithConfig func() string) error {
	envMap, err := godotenv.Unmarshal(stringWithConfig())
	if err != nil {
		return err
	}
	return godotenv.Write(envMap, "./config/config.env")

}
