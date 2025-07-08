package services

import (
	"strconv"
	"time"
	"xchat/config"
	"xchat/dto"
	"xchat/handlers"
	"xchat/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type authService struct {
	cfg    *config.Config
	secret []byte
	ur     UsersRepository
	jr     JwtRepository
}

// NewAuthService creates a new instance of authService with the provided secret.
func NewAuthService(cfg *config.Config, secret []byte, ur UsersRepository, jr JwtRepository) handlers.AuthService {
	return &authService{
		cfg:    cfg,
		secret: secret,
		ur:     ur,
		jr:     jr,
	}
}

func (as *authService) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(as.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (as *authService) Signup(req *dto.SignupReq) (at, rt string, err error) {
	exists, err := as.ur.UserExistsByUsername(req.Username)
	if err != nil {
		return "", "", err
	}
	if exists {
		return "", "", handlers.ErrUserAlreadyExists
	}

	exists, err = as.ur.UserExistsByEmail(req.Email)
	if err != nil {
		return "", "", err
	}
	if exists {
		return "", "", handlers.ErrEmailAlreadyExists
	}

	// check if password is strong
	if !isPasswordSecure(req.Password) {
		return "", "", handlers.ErrInsecurePassword
	}

	// hash password
	req.Password, err = hashPassword(req.Password)
	if err != nil {
		return "", "", err
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     models.UserRole, // Default role
	}

	err = as.ur.SaveUser(&user)
	if err != nil {
		return "", "", err
	}

	ATClaims := models.ATClaims{
		Username: user.Username,
		Email:    user.Email,
		Role:     models.UserRole,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(user.ID), // Use user ID as subject
			ExpiresAt: jwt.TimeFunc().Add(time.Duration(as.cfg.Jwt.TokenAge) * time.Minute).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "beta_service",
		},
	}
	at, err = as.GenerateToken(&ATClaims)
	if err != nil {
		return "", "", err
	}

	jti := uuid.NewString()
	RTClaims := models.RTClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        jti,
			Subject:   strconv.Itoa(user.ID), // Use user ID as subject
			ExpiresAt: jwt.TimeFunc().Add(time.Duration(as.cfg.Jwt.RefreshAge) * time.Minute).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "beta_service",
		},
	}
	rt, err = as.GenerateToken(&RTClaims)
	if err != nil {
		return "", "", err
	}

	// Save refresh token JTI to database
	refreshToken := &models.RefreshToken{
		JTI:       jti,
		UserID:    user.ID,
		ExpiresAt: time.Unix(RTClaims.ExpiresAt, 0),
		CreatedAt: time.Now(),
	}
	err = as.jr.SaveRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (as *authService) Login(req *dto.LoginReq) (at, rt string, err error) {
	user, err := as.ur.GetUserByUsername(req.Username)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", handlers.ErrUserDoesntExist
	}

	if !comparePassword(user.Password, req.Password) {
		return "", "", handlers.ErrIncorrectPassword
	}

	ATClaims := models.ATClaims{
		Username: user.Username,
		Email:    user.Email,
		Role:     models.UserRole,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(user.ID), // Use user ID as subject
			ExpiresAt: jwt.TimeFunc().Add(time.Duration(as.cfg.Jwt.TokenAge) * time.Minute).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "beta_service",
		},
	}
	at, err = as.GenerateToken(&ATClaims)
	if err != nil {
		return "", "", err
	}

	jti := uuid.NewString()
	RTClaims := models.RTClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        jti,
			Subject:   strconv.Itoa(user.ID), // Use user ID as subject
			ExpiresAt: jwt.TimeFunc().Add(time.Duration(as.cfg.Jwt.RefreshAge) * time.Minute).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "beta_service",
		},
	}
	rt, err = as.GenerateToken(&RTClaims)
	if err != nil {
		return "", "", err
	}

	// Save refresh token JTI to database
	refreshToken := &models.RefreshToken{
		JTI:       jti,
		UserID:    user.ID,
		ExpiresAt: time.Unix(RTClaims.ExpiresAt, 0),
		CreatedAt: time.Now(),
	}
	err = as.jr.SaveRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

// RefreshToken validates the refresh token and issues new tokens
func (as *authService) RefreshToken(refreshToken string) (at, rt string, err error) {
	// Parse the refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &models.RTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return as.secret, nil
	})
	if err != nil {
		return "", "", handlers.ErrInvalidToken
	}

	claims, ok := token.Claims.(*models.RTClaims)
	if !ok || !token.Valid {
		return "", "", handlers.ErrInvalidToken
	}

	// Check if the refresh token JTI exists in the database
	storedToken, err := as.jr.GetRefreshToken(claims.Id)
	if err != nil {
		return "", "", err
	}
	if storedToken == nil {
		return "", "", handlers.ErrInvalidToken
	}

	// Check if the token has expired
	if time.Now().After(storedToken.ExpiresAt) {
		// Clean up expired token
		as.jr.DeleteRefreshToken(claims.Id)
		return "", "", handlers.ErrTokenExpired
	}

	// Get user information
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return "", "", handlers.ErrInvalidToken
	}

	user, err := as.ur.GetUserByID(userID)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", handlers.ErrUserDoesntExist
	}

	// Delete the old refresh token (refresh token rotation)
	err = as.jr.DeleteRefreshToken(claims.Id)
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	ATClaims := models.ATClaims{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(user.ID),
			ExpiresAt: jwt.TimeFunc().Add(time.Duration(as.cfg.Jwt.TokenAge) * time.Minute).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "beta_service",
		},
	}
	at, err = as.GenerateToken(&ATClaims)
	if err != nil {
		return "", "", err
	}

	jti := uuid.NewString()
	RTClaims := models.RTClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        jti,
			Subject:   strconv.Itoa(user.ID),
			ExpiresAt: jwt.TimeFunc().Add(time.Duration(as.cfg.Jwt.RefreshAge) * time.Minute).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "beta_service",
		},
	}
	rt, err = as.GenerateToken(&RTClaims)
	if err != nil {
		return "", "", err
	}

	// Save new refresh token JTI to database
	newRefreshToken := &models.RefreshToken{
		JTI:       jti,
		UserID:    user.ID,
		ExpiresAt: time.Unix(RTClaims.ExpiresAt, 0),
		CreatedAt: time.Now(),
	}
	err = as.jr.SaveRefreshToken(newRefreshToken)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

// DeleteRefreshToken removes a refresh token from the database
func (as *authService) DeleteRefreshToken(jti string) error {
	return as.jr.DeleteRefreshToken(jti)
}
