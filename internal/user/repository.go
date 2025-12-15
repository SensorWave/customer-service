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
	curl -X POST http://localhost:8080/companies/<COMPANY_ID>/users \
	    -H "Content-Type: application/json" \
	    -d '{
	        "first_name": "John",
	        "last_name": "Doe",
	        "email": "john@acme.com",
	        "phone": "0611223344"
	    }'
*/
func (r *Repository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (company_id, first_name, last_name, email, phone, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
		RETURNING id
	`

	return r.DB.QueryRowContext(ctx, query,
		u.CompanyID, u.FirstName, u.LastName, u.Email, u.Phone,
	).Scan(&u.ID)
}

/* ---------- GET BY ID ---------- */

/*
curl http://localhost:8080/users/<USER_ID>
*/
func (r *Repository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `SELECT id, company_id, first_name, last_name, email, phone, created_at, updated_at
              FROM users WHERE id=$1`

	u := &User{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.CompanyID, &u.FirstName, &u.LastName, &u.Email, &u.Phone,
		&u.CreatedAt, &u.UpdatedAt,
	)

	return u, err
}

/* ---------- GET BY COMPANY ---------- */

/*
curl http://localhost:8080/companies/<COMPANY_ID>/users
*/
func (r *Repository) GetByCompany(ctx context.Context, companyID int) ([]User, error) {
	query := `SELECT id, company_id, first_name, last_name, email, phone, created_at, updated_at
              FROM users WHERE company_id=$1`

	rows, err := r.DB.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.CompanyID, &u.FirstName, &u.LastName, &u.Email, &u.Phone,
			&u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

/* ---------- UPDATE ---------- */

/*
curl -X PUT http://localhost:8080/users/<USER_ID> \
    -H "Content-Type: application/json" \
    -d '{
        "company_id": "3",
        "first_name": "Johnny",
        "last_name": "Doe",
        "email": "johnny.doe@acme.com",
        "phone": "0699887766"
    }'
*/

func (r *Repository) Update(ctx context.Context, u *User) error {
	query := `UPDATE users SET first_name=$2, last_name=$3, email=$4, phone=$5,
              updated_at=NOW() WHERE id=$1`

	res, err := r.DB.ExecContext(ctx, query,
		u.ID, u.FirstName, u.LastName, u.Email, u.Phone,
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
curl -X DELETE http://localhost:8080/users/<USER_ID>
*/
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
