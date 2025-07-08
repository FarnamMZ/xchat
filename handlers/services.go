package handlers

import (
	"xchat/dto"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	GenerateToken(claims jwt.Claims) (string, error)
	// Returning access token and a refresh token: at & rt strings.
	Signup(req *dto.SignupReq) (at, rt string, err error)
	Login(req *dto.LoginReq) (at, rt string, err error)
	RefreshToken(refreshToken string) (at, rt string, err error)
	DeleteRefreshToken(jti string) error
}
