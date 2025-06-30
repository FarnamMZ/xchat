package models

type User struct {
	Username string
	Password string
	Email    string
	Role     int
}

const (
	UserRole = iota
	AdminRole
)
