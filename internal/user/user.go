package user

import "time"

type User struct {
	ID        int       `json:"id"`
	CompanyID int       `json:"company_id"`
	IDAuthKC  string    `json:"id_auth_kc"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
