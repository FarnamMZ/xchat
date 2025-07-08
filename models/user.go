package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Don't expose password in JSON
	Email    string `json:"email"`
	Role     int    `json:"role"`
}

const (
	UserRole = iota
	AdminRole
)
