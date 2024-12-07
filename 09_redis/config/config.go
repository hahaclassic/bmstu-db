package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "./.env"
)

type Config struct {
	Postres PostgresConfig
	Redis   RedisConfig
}

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"`
}

type RedisConfig struct {
	Host       string        `env:"REDIS_HOST"`
	Port       string        `env:"REDIS_PORT"`
	Expiration time.Duration `env:"REDIS_EXPIRATION"`
}

func MustLoad() *Config {
	config := &Config{}

	err := cleanenv.ReadConfig(configPath, config)
	if err != nil {
		log.Fatalf("Error while loading config: %s", err)
	}

	return config
}
