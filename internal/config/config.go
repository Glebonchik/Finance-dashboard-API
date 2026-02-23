package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MLService MLServiceConfig
}

type ServerConfig struct {
	Port string `envconfig:"SERVER_PORT" default:":8080"`
	Env  string `envconfig:"SERVER_ENV" default:"development"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	DBName   string `envconfig:"DB_NAME" default:"finance_dashboard"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

// DSN возвращает строку подключения к БД
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD"`
}

// Address возвращает адрес Redis для подключения
func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type JWTConfig struct {
	Secret         string        `envconfig:"JWT_SECRET" required:"true"`
	AccessExpiry   time.Duration `envconfig:"JWT_ACCESS_EXPIRATION" default:"15m"`
	RefreshExpiry  time.Duration `envconfig:"JWT_REFRESH_EXPIRATION" default:"24h"`
}

type MLServiceConfig struct {
	Host string `envconfig:"ML_SERVICE_HOST" default:"localhost"`
	Port string `envconfig:"ML_SERVICE_PORT" default:"50051"`
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
