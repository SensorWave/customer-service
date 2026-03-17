package config

import (
	"strings"
	"testing"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()

	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_USER", "user")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "customer_service_db")
	t.Setenv("RABBITMQ_HOST", "rabbitmq")
}

func TestLoadConfig_Success(t *testing.T) {
	setRequiredEnv(t)

	_, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
}

func TestLoadConfig_MissingRequiredEnv(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("RABBITMQ_HOST", "")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected LoadConfig to fail when a required env var is missing")
	}

	if !strings.Contains(err.Error(), "RABBITMQ_HOST") {
		t.Fatalf("expected error to mention RABBITMQ_HOST, got %v", err)
	}
}
