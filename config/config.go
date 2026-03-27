package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	RabbitMQ   string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		RabbitMQ:   os.Getenv("RABBITMQ_HOST"),
	}

	var missing []string
	requiredValues := []struct {
		name  string
		value string
	}{
		{name: "DB_HOST", value: cfg.DBHost},
		{name: "DB_USER", value: cfg.DBUser},
		{name: "DB_PASSWORD", value: cfg.DBPassword},
		{name: "DB_NAME", value: cfg.DBName},
		{name: "RABBITMQ_HOST", value: cfg.RabbitMQ},
	}

	for _, item := range requiredValues {
		if item.value == "" {
			missing = append(missing, item.name)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) RabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:5672/", c.RabbitMQ)
}
