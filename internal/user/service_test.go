package user

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	company "customer-service/internal/company"

	"github.com/DATA-DOG/go-sqlmock"
)

type publishedEvent struct {
	routingKey string
	payload    interface{}
}

type spyPublisher struct {
	events []publishedEvent
}

func (s *spyPublisher) Publish(routingKey string, payload interface{}) {
	s.events = append(s.events, publishedEvent{routingKey: routingKey, payload: payload})
}

func newTestService(t *testing.T) (*Service, sqlmock.Sqlmock, *spyPublisher, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	spy := &spyPublisher{}
	svc := NewService(
		NewRepository(db),
		company.NewRepository(db),
		spy,
	)

	cleanup := func() {
		_ = db.Close()
	}

	return svc, mock, spy, cleanup
}

func TestService_CreateUser_Success(t *testing.T) {
	svc, mock, spy, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	u := &User{
		CompanyID: 7,
		IDAuthKC:  "kc-acme-001",
		Role:      "admin",
	}

	mock.ExpectQuery(`SELECT id, name, email, phone, address, created_at, updated_at\s+FROM companies\s+WHERE id = \$1`).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "email", "phone", "address", "created_at", "updated_at",
		}).AddRow("7", "Acme", "info@acme.com", "0600000000", "Main St", now, now))

	mock.ExpectQuery(`INSERT INTO user_companies \(company_id, id_auth_kc, role, created_at, updated_at\)\s+VALUES \(\$1,\$2,\$3,NOW\(\),NOW\(\)\)\s+RETURNING id`).
		WithArgs(u.CompanyID, u.IDAuthKC, u.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(12))

	err := svc.CreateUser(context.Background(), u)
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	if u.ID != 12 {
		t.Fatalf("expected user ID to be set to 12, got %d", u.ID)
	}
	if len(spy.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(spy.events))
	}
	if spy.events[0].routingKey != "UserCreated" {
		t.Fatalf("expected UserCreated event, got %s", spy.events[0].routingKey)
	}
	if spy.events[0].payload != u {
		t.Fatal("expected published payload to be the created user pointer")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_CreateUser_CompanyLookupError(t *testing.T) {
	svc, mock, spy, cleanup := newTestService(t)
	defer cleanup()

	u := &User{CompanyID: 99}

	mock.ExpectQuery(`SELECT id, name, email, phone, address, created_at, updated_at\s+FROM companies\s+WHERE id = \$1`).
		WithArgs(99).
		WillReturnError(errors.New("db down"))

	err := svc.CreateUser(context.Background(), u)
	if err == nil {
		t.Fatal("expected error from CreateUser, got nil")
	}
	if len(spy.events) != 0 {
		t.Fatalf("expected no published events, got %d", len(spy.events))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetUserByID(t *testing.T) {
	svc, mock, _, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, company_id, id_auth_kc, role, created_at, updated_at\s+FROM user_companies WHERE id=\$1`).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "company_id", "id_auth_kc", "role", "created_at", "updated_at",
		}).AddRow(5, 7, "kc-acme-001", "admin", now, now))

	u, err := svc.GetUserByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}
	if u.ID != 5 || u.CompanyID != 7 || u.IDAuthKC != "kc-acme-001" || u.Role != "admin" {
		t.Fatalf("unexpected user returned: %#v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetUserByKeycloakID(t *testing.T) {
	svc, mock, _, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, company_id, id_auth_kc, role, created_at, updated_at\s+FROM user_companies WHERE id_auth_kc=\$1`).
		WithArgs("kc-acme-001").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "company_id", "id_auth_kc", "role", "created_at", "updated_at",
		}).AddRow(5, 7, "kc-acme-001", "admin", now, now))

	u, err := svc.GetUserByKeycloakID(context.Background(), "kc-acme-001")
	if err != nil {
		t.Fatalf("GetUserByKeycloakID returned error: %v", err)
	}
	if u.ID != 5 || u.CompanyID != 7 || u.IDAuthKC != "kc-acme-001" || u.Role != "admin" {
		t.Fatalf("unexpected user returned: %#v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetUserRoleByKeycloakID(t *testing.T) {
	svc, mock, _, cleanup := newTestService(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT role FROM user_companies WHERE id_auth_kc=\$1`).
		WithArgs("kc-acme-001").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("manager"))

	role, err := svc.GetUserRoleByKeycloakID(context.Background(), "kc-acme-001")
	if err != nil {
		t.Fatalf("GetUserRoleByKeycloakID returned error: %v", err)
	}
	if role != "manager" {
		t.Fatalf("expected role manager, got %s", role)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetUserCompanyIDByKeycloakID(t *testing.T) {
	svc, mock, _, cleanup := newTestService(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT company_id FROM user_companies WHERE id_auth_kc=\$1`).
		WithArgs("kc-acme-001").
		WillReturnRows(sqlmock.NewRows([]string{"company_id"}).AddRow(42))

	companyID, err := svc.GetUserCompanyIDByKeycloakID(context.Background(), "kc-acme-001")
	if err != nil {
		t.Fatalf("GetUserCompanyIDByKeycloakID returned error: %v", err)
	}
	if companyID != 42 {
		t.Fatalf("expected company ID 42, got %d", companyID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetUsersByCompany(t *testing.T) {
	svc, mock, _, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, company_id, id_auth_kc, role, created_at, updated_at\s+FROM user_companies WHERE company_id=\$1`).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "company_id", "id_auth_kc", "role", "created_at", "updated_at",
		}).
			AddRow(1, 7, "kc-acme-001", "admin", now, now).
			AddRow(2, 7, "kc-acme-002", "user", now, now))

	users, err := svc.GetUsersByCompany(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetUsersByCompany returned error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_UpdateUser_Success(t *testing.T) {
	svc, mock, spy, cleanup := newTestService(t)
	defer cleanup()

	u := &User{
		ID:       10,
		IDAuthKC: "kc-acme-009",
		Role:     "manager",
	}

	mock.ExpectExec(`UPDATE user_companies SET id_auth_kc=\$2, role=\$3,\s+updated_at=NOW\(\) WHERE id=\$1`).
		WithArgs(u.ID, u.IDAuthKC, u.Role).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := svc.UpdateUser(context.Background(), u)
	if err != nil {
		t.Fatalf("UpdateUser returned error: %v", err)
	}
	if len(spy.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(spy.events))
	}
	if spy.events[0].routingKey != "UserUpdated" {
		t.Fatalf("expected UserUpdated event, got %s", spy.events[0].routingKey)
	}
	if spy.events[0].payload != u {
		t.Fatal("expected published payload to be the updated user pointer")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_UpdateUser_NotFound(t *testing.T) {
	svc, mock, spy, cleanup := newTestService(t)
	defer cleanup()

	u := &User{ID: 123}

	mock.ExpectExec(`UPDATE user_companies SET id_auth_kc=\$2, role=\$3,\s+updated_at=NOW\(\) WHERE id=\$1`).
		WithArgs(u.ID, u.IDAuthKC, u.Role).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := svc.UpdateUser(context.Background(), u)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
	if len(spy.events) != 0 {
		t.Fatalf("expected no published events, got %d", len(spy.events))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_DeleteUser(t *testing.T) {
	svc, mock, spy, cleanup := newTestService(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM user_companies WHERE id=\$1`).
		WithArgs(77).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := svc.DeleteUser(context.Background(), 77)
	if err != nil {
		t.Fatalf("DeleteUser returned error: %v", err)
	}
	if len(spy.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(spy.events))
	}
	if spy.events[0].routingKey != "UserDeleted" {
		t.Fatalf("expected UserDeleted event, got %s", spy.events[0].routingKey)
	}
	expectedPayload := map[string]string{"id": "77"}
	if !reflect.DeepEqual(spy.events[0].payload, expectedPayload) {
		t.Fatalf("unexpected delete payload: %#v", spy.events[0].payload)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_DeleteUser_Error(t *testing.T) {
	svc, mock, spy, cleanup := newTestService(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM user_companies WHERE id=\$1`).
		WithArgs(77).
		WillReturnError(errors.New("delete failed"))

	err := svc.DeleteUser(context.Background(), 77)
	if err == nil {
		t.Fatal("expected error from DeleteUser, got nil")
	}
	if len(spy.events) != 0 {
		t.Fatalf("expected no published events, got %d", len(spy.events))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
