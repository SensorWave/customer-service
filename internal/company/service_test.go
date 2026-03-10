package company

import (
	"context"
	"errors"
	"testing"
	"time"

	"customer-service/internal/event"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestService(t *testing.T) (*Service, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	svc := NewService(NewRepository(db), &event.Publisher{})

	cleanup := func() {
		_ = db.Close()
	}

	return svc, mock, cleanup
}

func TestService_CreateCompany_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	c := &Company{
		Name:    "Acme Corp",
		Email:   "contact@acme.com",
		Phone:   "0601020304",
		Address: "42 rue de Paris",
	}

	mock.ExpectQuery(`INSERT INTO companies \(name, email, phone, address, created_at, updated_at\)\s+VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\)\s+RETURNING id`).
		WithArgs(c.Name, c.Email, c.Phone, c.Address).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("3"))

	err := svc.CreateCompany(context.Background(), c)
	if err != nil {
		t.Fatalf("CreateCompany returned error: %v", err)
	}
	if c.ID != "3" {
		t.Fatalf("expected company ID to be set to 3, got %s", c.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_CreateCompany_Error(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	c := &Company{Name: "Acme Corp"}

	mock.ExpectQuery(`INSERT INTO companies \(name, email, phone, address, created_at, updated_at\)\s+VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\)\s+RETURNING id`).
		WithArgs(c.Name, c.Email, c.Phone, c.Address).
		WillReturnError(errors.New("insert failed"))

	err := svc.CreateCompany(context.Background(), c)
	if err == nil {
		t.Fatal("expected error from CreateCompany, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetCompanyByID(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, name, email, phone, address, created_at, updated_at\s+FROM companies\s+WHERE id = \$1`).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "email", "phone", "address", "created_at", "updated_at",
		}).AddRow("5", "Acme Corp", "contact@acme.com", "0601020304", "42 rue de Paris", now, now))

	c, err := svc.GetCompanyByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetCompanyByID returned error: %v", err)
	}
	if c == nil || c.ID != "5" {
		t.Fatalf("unexpected company returned: %#v", c)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_GetAllCompanies(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, name, email, phone, address, created_at, updated_at\s+FROM companies\s+ORDER BY created_at DESC`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "email", "phone", "address", "created_at", "updated_at",
		}).
			AddRow("1", "A", "a@acme.com", "0600000001", "Addr A", now, now).
			AddRow("2", "B", "b@acme.com", "0600000002", "Addr B", now, now))

	companies, err := svc.GetAllCompanies(context.Background())
	if err != nil {
		t.Fatalf("GetAllCompanies returned error: %v", err)
	}
	if len(companies) != 2 {
		t.Fatalf("expected 2 companies, got %d", len(companies))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_UpdateCompany_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	c := &Company{
		ID:      "3",
		Name:    "Acme Corporation Updated",
		Email:   "support@acme.com",
		Phone:   "0707070707",
		Address: "100 avenue de Lyon",
	}

	mock.ExpectExec(`UPDATE companies\s+SET name=\$1, email=\$2, phone=\$3, address=\$4, updated_at=NOW\(\)\s+WHERE id=\$5`).
		WithArgs(c.Name, c.Email, c.Phone, c.Address, c.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := svc.UpdateCompany(context.Background(), c)
	if err != nil {
		t.Fatalf("UpdateCompany returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_UpdateCompany_NotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	c := &Company{ID: "999"}

	mock.ExpectExec(`UPDATE companies\s+SET name=\$1, email=\$2, phone=\$3, address=\$4, updated_at=NOW\(\)\s+WHERE id=\$5`).
		WithArgs(c.Name, c.Email, c.Phone, c.Address, c.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := svc.UpdateCompany(context.Background(), c)
	if err == nil || err.Error() != "company not found" {
		t.Fatalf("expected company not found error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_DeleteCompany_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM companies WHERE id=\$1`).
		WithArgs(77).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := svc.DeleteCompany(context.Background(), 77)
	if err != nil {
		t.Fatalf("DeleteCompany returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestService_DeleteCompany_NotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM companies WHERE id=\$1`).
		WithArgs(77).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := svc.DeleteCompany(context.Background(), 77)
	if err == nil || err.Error() != "company not found" {
		t.Fatalf("expected company not found error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
