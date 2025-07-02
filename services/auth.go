package services

import (
	"github.com/golang-jwt/jwt"
	"time"
	"xchat/dto"
	"xchat/handlers"
	"xchat/models"
)

type authService struct {
	secret []byte
	ur     UsersRepository
}

// NewAuthService creates a new instance of authService with the provided secret.
func NewAuthService(secret []byte, ur UsersRepository) handlers.AuthService {
	return &authService{
		secret: secret,
		ur:     ur,
	}
}

func (as *authService) GenerateToken(claims *models.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(as.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (as *authService) Signup(req *dto.SignupReq) (string, error) {
	exists, err := as.ur.UserExistsByUsername(req.Username)
	if err != nil {
		return "", err
	}
	if exists {
		return "", handlers.ErrUserAlreadyExists
	}

	exists, err = as.ur.UserExistsByEmail(req.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", handlers.ErrEmailAlreadyExists
	}

	// check if password is strong
	if !isPasswordSecure(req.Password) {
		return "", handlers.ErrInsecurePassword
	}

	// hash password
	req.Password, err = hashPassword(req.Password)
	if err != nil {
		return "", err
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     models.UserRole, // Default role
	}

	err = as.ur.SaveUser(&user)
	if err != nil {
		return "", err
	}

	claims := models.Claims{
		Username: user.Username,
		Email:    user.Email,
		Role:     models.UserRole,
		StandardClaims: jwt.StandardClaims{
			// 1 hour expiration
			ExpiresAt: jwt.TimeFunc().Add(1 * time.Hour).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "beta_service",
		},
	}
	token, err := as.GenerateToken(&claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (as *authService) Login(req *dto.LoginReq) (string, error) {
	user, err := as.ur.GetUserByUsername(req.Username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", handlers.ErrUserDoesntExist
	}

	if !comparePassword(user.Password, req.Password) {
		return "", handlers.ErrIncorrectPassword
	}

	claims := models.Claims{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			// 1 hour expiration
			ExpiresAt: jwt.TimeFunc().Add(1 * time.Hour).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "beta_service",
		},
	}
	token, err := as.GenerateToken(&claims)
	if err != nil {
		return "", err
	}

	return token, nil
}
