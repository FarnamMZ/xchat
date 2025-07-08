package models

import "github.com/golang-jwt/jwt"

// ATClaims represents the JWT claims used for authentication.
type ATClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}

// RTClaims represents the JWT claims used for refresh token.
// only containing jwt.StandardClaims and is basically empty.
type RTClaims struct {
	jwt.StandardClaims
}
