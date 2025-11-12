package user

import "database/sql"

type Repository struct {
    DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{DB: db}
}

func (r *Repository) Create(u *User) error {
    query := `INSERT INTO users (id, company_id, first_name, last_name, email, phone, created_at, updated_at)
              VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())`
    _, err := r.DB.Exec(query, u.ID, u.CompanyID, u.FirstName, u.LastName, u.Email, u.Phone)
    return err
}

func (r *Repository) GetByID(id string) (*User, error) {
    u := &User{}
    query := `SELECT id, company_id, first_name, last_name, email, phone, created_at, updated_at FROM users WHERE id=$1`
    err := r.DB.QueryRow(query, id).Scan(&u.ID, &u.CompanyID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.CreatedAt, &u.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return u, nil
}
