package config

import (
	"log/slog"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Note: OTLP is configured using its internal configuration mechanisms

type ServerConfig struct {
	ListenAddress string   `yaml:"listen_address" env:"LISTEN_ADDRESS" env-default:"0.0.0.0:8080"`
	CorsOrigins   []string `yaml:"cors_origins" env:"CORS_ORIGINS" env-default:"http://localhost:3000,http://localhost:8080,http://localhost,https://store-it.ru,https://www.store-it.ru,http://store-it.ru,http://www.store-it.ru"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"storeit"`
	User     string `yaml:"user" env:"DB_USER" env-default:"storeit"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"storeit"`
}

type YandexOAuthConfig struct {
	ClientID     string `yaml:"client_id" env:"YANDEX_OAUTH_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" env:"YANDEX_OAUTH_CLIENT_SECRET"`
}

type KafkaConfig struct {
	// Comma-separated list of brokers, e.g. "localhost:9092,localhost:9093"
	Brokers string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	// Whether Kafka integration is enabled
	Enabled bool `yaml:"enabled" env:"KAFKA_ENABLED" env-default:"false"`
	// Topic for audit events
	AuditTopic string `yaml:"audit_topic" env:"KAFKA_AUDIT_TOPIC" env-default:"audit.object-changes"`
}

// GetBrokersList returns the list of Kafka brokers
func (k *KafkaConfig) GetBrokersList() []string {
	return strings.Split(k.Brokers, ",")
}

type Config struct {
	ServiceName string            `yaml:"service_name" env:"SERVICE_NAME" env-default:"storeit-backend"`
	Server      ServerConfig      `yaml:"server"`
	Database    DatabaseConfig    `yaml:"database"`
	YandexOAuth YandexOAuthConfig `yaml:"yandex_oauth"`
	Kafka       KafkaConfig       `yaml:"kafka"`
}

func GetConfigOrDie() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		slog.Error("Failed to read config", "error", err)
	}

	return &cfg
}
