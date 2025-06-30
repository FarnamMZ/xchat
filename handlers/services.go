package handlers

import (
	"beta/dto"
	"beta/models"
)

type AuthService interface {
	GenerateToken(claims models.Claims) (string, error)
	Signup(req dto.SignupReq) (string, error)
}
