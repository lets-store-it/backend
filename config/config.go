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

type TelemetryConfig struct {
	Tracing TracingConfig `yaml:"tracing"`
}

type TracingConfig struct {
	Endpoint string `yaml:"endpoint" env:"OTEL_EXPORTER_OTLP_ENDPOINT" env-default:"otlp-gateway-prod-eu-west-2.grafana.net/otlp"`
	Headers  string `yaml:"headers" env:"OTEL_EXPORTER_OTLP_HEADERS" env-default:"Authorization=Basic Nzg3MDQ0OmdsY19leUp2SWpvaU9UZzNPVGMzSWl3aWJpSTZJbk4wWVdOckxUYzROekEwTkMxdmRHeHdMWGR5YVhSbExXUmxkaUlzSW1zaU9pSkNVSFJGT0RjMGRXSkRiVEpNUlRGdWJVWTRORFUyTUZRaUxDSnRJanA3SW5JaU9pSndjbTlrTFdWMUxYZGxjM1F0TWlKOWZRPT0="`
	Insecure bool   `yaml:"insecure" env:"OTEL_EXPORTER_OTLP_INSECURE" env-default:"false"`
}

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Database    DatabaseConfig    `yaml:"database"`
	YandexOAuth YandexOAuthConfig `yaml:"yandex_oauth"`
	Telemetry   TelemetryConfig   `yaml:"telemetry"`
}

func GetConfigOrDie() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
