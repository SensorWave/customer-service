package company

import "database/sql"

type Repository struct {
    DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{DB: db}
}

func (r *Repository) Create(c *Company) error {
    query := `INSERT INTO companies (id, name, email, phone, address, created_at, updated_at)
              VALUES ($1,$2,$3,$4,$5,NOW(),NOW())`
    _, err := r.DB.Exec(query, c.ID, c.Name, c.Email, c.Phone, c.Address)
    return err
}

func (r *Repository) GetByID(id string) (*Company, error) {
    c := &Company{}
    query := `SELECT id, name, email, phone, address, created_at, updated_at FROM companies WHERE id=$1`
    err := r.DB.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.CreatedAt, &c.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return c, nil
}
