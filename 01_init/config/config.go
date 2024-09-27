package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "./.env"
)

type Config struct {
	Postres   PostgresConfig
	Generator GeneratorConfig
}

type GeneratorConfig struct {
	OutputCSV       string // path/to/folder
	DeleteCmd       bool
	RecordsPerTable int
}

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"`
}

func MustLoad() *Config {
	config := &Config{}

	err := cleanenv.ReadConfig(configPath, config)
	if err != nil {
		log.Fatalf("Error while loading config: %s", err)
	}

	return config
}
