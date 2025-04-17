package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	ListenAddress string `yaml:"listen_address" env:"LISTEN_ADDRESS" env-default:"0.0.0.0:8080"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"storeit"`
	User     string `yaml:"user" env:"DB_USER" env-default:"storeit"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"storeit"`
}

type YandexOAuthConfig struct {
	ClientID     string `yaml:"client_id" env:"YANDEX_OAUTH_CLIENT_ID" env-default:"712925a705b34f5399ba6f067347266b"`
	ClientSecret string `yaml:"client_secret" env:"YANDEX_OAUTH_CLIENT_SECRET" env-default:"331b7a0292044e958989397a28c56bc7"`
}

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	YandexOAuth YandexOAuthConfig
}

func GetConfigOrDie() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
