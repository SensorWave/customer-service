package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	company "customer-service/internal/company"
	"customer-service/internal/event"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestService(t *testing.T) (*Service, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	svc := NewService(
		NewRepository(db),
		company.NewRepository(db),
		&event.Publisher{},
	)

	cleanup := func() {
		_ = db.Close()
	}

	return svc, mock, cleanup
}

func TestService_CreateUser_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	u := &User{
		CompanyID: 7,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@acme.com",
		Phone:     "0611223344",
	}

	mock.ExpectQuery(`SELECT id, name, email, phone, address, created_at, updated_at\s+FROM companies\s+WHERE id = \$1`).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "email", "phone", "address", "created_at", "updated_at",
		}).AddRow("7", "Acme", "info@acme.com", "0600000000", "Main St", now, now))

	mock.ExpectQuery(`INSERT INTO users \(company_id, first_name, last_name, email, phone, created_at, updated_at\)\s+VALUES \(\$1,\$2,\$3,\$4,\$5,NOW\(\),NOW\(\)\)\s+RETURNING id`).
		WithArgs(u.CompanyID, u.FirstName, u.LastName, u.Email, u.Phone).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(12))

	err := svc.CreateUser(context.Background(), u)
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	if u.ID != 12 {
		t.Fatalf("expected user ID to be set to 12, got %d", u.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_CreateUser_CompanyLookupError(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	u := &User{CompanyID: 99}

	mock.ExpectQuery(`SELECT id, name, email, phone, address, created_at, updated_at\s+FROM companies\s+WHERE id = \$1`).
		WithArgs(99).
		WillReturnError(errors.New("db down"))

	err := svc.CreateUser(context.Background(), u)
	if err == nil {
		t.Fatal("expected error from CreateUser, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetUserByID(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, company_id, first_name, last_name, email, phone, created_at, updated_at\s+FROM users WHERE id=\$1`).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "company_id", "first_name", "last_name", "email", "phone", "created_at", "updated_at",
		}).AddRow(5, 7, "Jane", "Doe", "jane@acme.com", "0611223344", now, now))

	u, err := svc.GetUserByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}
	if u.ID != 5 || u.CompanyID != 7 {
		t.Fatalf("unexpected user returned: %#v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetUsersByCompany(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, company_id, first_name, last_name, email, phone, created_at, updated_at\s+FROM users WHERE company_id=\$1`).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "company_id", "first_name", "last_name", "email", "phone", "created_at", "updated_at",
		}).
			AddRow(1, 7, "A", "One", "a@acme.com", "0600000001", now, now).
			AddRow(2, 7, "B", "Two", "b@acme.com", "0600000002", now, now))

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
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	u := &User{
		ID:        10,
		FirstName: "Updated",
		LastName:  "Name",
		Email:     "updated@acme.com",
		Phone:     "0699887766",
	}

	mock.ExpectExec(`UPDATE users SET first_name=\$2, last_name=\$3, email=\$4, phone=\$5,\s+updated_at=NOW\(\) WHERE id=\$1`).
		WithArgs(u.ID, u.FirstName, u.LastName, u.Email, u.Phone).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := svc.UpdateUser(context.Background(), u)
	if err != nil {
		t.Fatalf("UpdateUser returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_UpdateUser_NotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	u := &User{ID: 123}

	mock.ExpectExec(`UPDATE users SET first_name=\$2, last_name=\$3, email=\$4, phone=\$5,\s+updated_at=NOW\(\) WHERE id=\$1`).
		WithArgs(u.ID, u.FirstName, u.LastName, u.Email, u.Phone).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := svc.UpdateUser(context.Background(), u)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_DeleteUser(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM users WHERE id=\$1`).
		WithArgs(77).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := svc.DeleteUser(context.Background(), 77)
	if err != nil {
		t.Fatalf("DeleteUser returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_DeleteUser_Error(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM users WHERE id=\$1`).
		WithArgs(77).
		WillReturnError(errors.New("delete failed"))

	err := svc.DeleteUser(context.Background(), 77)
	if err == nil {
		t.Fatal("expected error from DeleteUser, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
