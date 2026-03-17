package user

import (
	"context"
	"database/sql"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

/* ---------- CREATE ---------- */

/*
	curl -X POST http://localhost:8080/companies/<COMPANY_ID>/users-company \
	    -H "Content-Type: application/json" \
	    -d '{
	        "id_auth_kc": "kc-user-001",
	        "role": "admin"
	    }'
*/
func (r *Repository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO user_companies (company_id, id_auth_kc, role, created_at, updated_at)
		VALUES ($1,$2,$3,NOW(),NOW())
		RETURNING id
	`

	return r.DB.QueryRowContext(ctx, query,
		u.CompanyID, u.IDAuthKC, u.Role,
	).Scan(&u.ID)
}

/* ---------- GET BY ID ---------- */

/*
curl http://localhost:8080/users-company/<USER_ID>
*/
func (r *Repository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `SELECT id, company_id, id_auth_kc, role, created_at, updated_at
              FROM user_companies WHERE id=$1`

	u := &User{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.CompanyID, &u.IDAuthKC, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)

	return u, err
}

/* ---------- GET BY COMPANY ---------- */

/*
curl http://localhost:8080/companies/<COMPANY_ID>/users-company
*/
func (r *Repository) GetByCompany(ctx context.Context, companyID int) ([]User, error) {
	query := `SELECT id, company_id, id_auth_kc, role, created_at, updated_at
              FROM user_companies WHERE company_id=$1`

	rows, err := r.DB.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.CompanyID, &u.IDAuthKC, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

/* ---------- UPDATE ---------- */

/*
curl -X PUT http://localhost:8080/users-company/<USER_ID> \
    -H "Content-Type: application/json" \
    -d '{
        "id_auth_kc": "kc-user-002",
        "role": "manager"
    }'
*/

func (r *Repository) Update(ctx context.Context, u *User) error {
	query := `UPDATE user_companies SET id_auth_kc=$2, role=$3,
              updated_at=NOW() WHERE id=$1`

	res, err := r.DB.ExecContext(ctx, query,
		u.ID, u.IDAuthKC, u.Role,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

/* ---------- DELETE ---------- */

/*
curl -X DELETE http://localhost:8080/users-company/<USER_ID>
*/
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM user_companies WHERE id=$1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
