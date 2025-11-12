package config

import (
    "fmt"
    "log"
    "os"
)

type Config struct {
    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    RabbitMQ   string
}

func LoadConfig() *Config {
    cfg := &Config{
        DBHost: os.Getenv("DB_HOST"),
        DBUser: os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName: os.Getenv("DB_NAME"),
        RabbitMQ: os.Getenv("RABBITMQ_HOST"),
    }

    if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.RabbitMQ == "" {
        log.Fatal("Missing required environment variables")
    }

    return cfg
}

func (c *Config) PostgresDSN() string {
    return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) RabbitMQURL() string {
    return fmt.Sprintf("amqp://%s:5672/", c.RabbitMQ)
}
