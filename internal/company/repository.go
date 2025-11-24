package company

import (
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// ---------- CREATE ----------

/*
	curl -X POST http://localhost:8080/companies \
	    -H "Content-Type: application/json" \
	    -d '{
	        "name": "Acme Corp",
	        "email": "contact@acme.com",
	        "phone": "0601020304",
	        "address": "42 rue de Paris"
	    }'
*/
func (r *Repository) Create(ctx context.Context, c *Company) error {
	query := `
        INSERT INTO companies (name, email, phone, address, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id
    `
	return r.DB.QueryRowContext(ctx, query,
		c.Name, c.Email, c.Phone, c.Address,
	).Scan(&c.ID)
}

// ---------- GET BY ID ----------

/*
curl http://localhost:8080/companies/<COMPANY_ID>
*/
func (r *Repository) GetByID(ctx context.Context, id int) (*Company, error) {
	query := `
        SELECT id, name, email, phone, address, created_at, updated_at
        FROM companies
        WHERE id = $1
    `

	c := &Company{}

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Email, &c.Phone,
		&c.Address, &c.CreatedAt, &c.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return c, err
}

// ---------- GET ALL ----------

/*
curl http://localhost:8080/companies
*/
func (r *Repository) GetAll(ctx context.Context) ([]*Company, error) {
	query := `
        SELECT id, name, email, phone, address, created_at, updated_at
        FROM companies
        ORDER BY created_at DESC
    `

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*Company{}

	for rows.Next() {
		c := &Company{}
		err := rows.Scan(
			&c.ID, &c.Name, &c.Email, &c.Phone,
			&c.Address, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	return companies, nil
}

// ---------- UPDATE ----------

/*
curl -X PUT http://localhost:8080/companies/<COMPANY_ID> \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Acme Corporation Updated",
        "email": "support@acme.com",
        "phone": "0707070707",
        "address": "100 avenue de Lyon"
    }'
*/

func (r *Repository) Update(ctx context.Context, c *Company) error {
	query := `
        UPDATE companies
        SET name=$1, email=$2, phone=$3, address=$4, updated_at=NOW()
        WHERE id=$5
    `
	res, err := r.DB.ExecContext(ctx, query,
		c.Name, c.Email, c.Phone, c.Address, c.ID,
	)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("company not found")
	}

	return nil
}

// ---------- DELETE ----------

/*
curl -X DELETE http://localhost:8080/companies/<COMPANY_ID>
*/
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM companies WHERE id=$1`

	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("company not found")
	}

	return nil
}
