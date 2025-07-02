package handlers

import (
	"xchat/dto"
	"xchat/models"
)

type AuthService interface {
	GenerateToken(claims *models.Claims) (string, error)
	Signup(req *dto.SignupReq) (string, error)
	Login(req *dto.LoginReq) (string, error)
}
