package models

import "github.com/golang-jwt/jwt"

// Claims represents the JWT claims used for authentication.
type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}
