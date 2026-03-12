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
	t.Setenv("KEYCLOAK_ISSUER", "http://localhost:8888/realms/dev")
	t.Setenv("KEYCLOAK_JWKS_URL", "http://localhost:8888/realms/dev/protocol/openid-connect/certs")
}

func TestLoadConfig_Success(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("KEYCLOAK_AUDIENCE", "customer-service")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.KeycloakIssuer != "http://localhost:8888/realms/dev" {
		t.Fatalf("unexpected Keycloak issuer: %q", cfg.KeycloakIssuer)
	}
	if cfg.KeycloakJWKSURL != "http://localhost:8888/realms/dev/protocol/openid-connect/certs" {
		t.Fatalf("unexpected Keycloak JWKS URL: %q", cfg.KeycloakJWKSURL)
	}
	if cfg.KeycloakAudience != "customer-service" {
		t.Fatalf("unexpected Keycloak audience: %q", cfg.KeycloakAudience)
	}
}

func TestLoadConfig_AllowsEmptyAudience(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("KEYCLOAK_AUDIENCE", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.KeycloakAudience != "" {
		t.Fatalf("expected empty audience, got %q", cfg.KeycloakAudience)
	}
}

func TestLoadConfig_MissingRequiredEnv(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("KEYCLOAK_JWKS_URL", "")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected LoadConfig to fail when a required env var is missing")
	}

	if !strings.Contains(err.Error(), "KEYCLOAK_JWKS_URL") {
		t.Fatalf("expected error to mention KEYCLOAK_JWKS_URL, got %v", err)
	}
}
