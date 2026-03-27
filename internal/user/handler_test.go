package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func newUserTestRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	userHandler := NewHandler(NewService(NewRepository(db), nil, nil))

	r := gin.New()
	r.GET("/users-company/keycloak/:keycloak_id", userHandler.GetUserByKeycloakID)
	r.GET("/users-company/keycloak/:keycloak_id/role", userHandler.GetUserRoleByKeycloakID)
	r.GET("/users-company/keycloak/:keycloak_id/company-id", userHandler.GetUserCompanyIDByKeycloakID)

	return r
}

func TestRoute_GetUserByKeycloakID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, company_id, id_auth_kc, role, created_at, updated_at\s+FROM user_companies WHERE id_auth_kc=\$1`).
		WithArgs("kc-acme-001").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "company_id", "id_auth_kc", "role", "created_at", "updated_at",
		}).AddRow(5, 7, "kc-acme-001", "admin", now, now))

	req := httptest.NewRequest(http.MethodGet, "/users-company/keycloak/kc-acme-001", nil)
	rec := httptest.NewRecorder()

	newUserTestRouter(db).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", rec.Code, rec.Body.String())
	}

	var got User
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.ID != 5 || got.CompanyID != 7 || got.IDAuthKC != "kc-acme-001" || got.Role != "admin" {
		t.Fatalf("unexpected response: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestRoute_GetUserRoleByKeycloakID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT role FROM user_companies WHERE id_auth_kc=\$1`).
		WithArgs("kc-acme-001").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("manager"))

	req := httptest.NewRequest(http.MethodGet, "/users-company/keycloak/kc-acme-001/role", nil)
	rec := httptest.NewRecorder()

	newUserTestRouter(db).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", rec.Code, rec.Body.String())
	}

	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got["role"] != "manager" {
		t.Fatalf("expected role manager, got %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestRoute_GetUserCompanyIDByKeycloakID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT company_id FROM user_companies WHERE id_auth_kc=\$1`).
		WithArgs("kc-acme-001").
		WillReturnRows(sqlmock.NewRows([]string{"company_id"}).AddRow(42))

	req := httptest.NewRequest(http.MethodGet, "/users-company/keycloak/kc-acme-001/company-id", nil)
	rec := httptest.NewRecorder()

	newUserTestRouter(db).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", rec.Code, rec.Body.String())
	}

	var got map[string]int
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got["company_id"] != 42 {
		t.Fatalf("expected company_id 42, got %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestRoute_GetUserRoleByKeycloakID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT role FROM user_companies WHERE id_auth_kc=\$1`).
		WithArgs("missing-user").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/users-company/keycloak/missing-user/role", nil)
	rec := httptest.NewRecorder()

	newUserTestRouter(db).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d with body %s", rec.Code, rec.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
