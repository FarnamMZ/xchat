package services

import "xchat/models"

type UsersRepository interface {
	UserExistsByUsername(username string) (bool, error)
	UserExistsByEmail(email string) (bool, error)
	SaveUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}

type JwtRepository interface {
	SaveRefreshToken(token *models.RefreshToken) error
	GetRefreshToken(jti string) (*models.RefreshToken, error)
	DeleteRefreshToken(jti string) error
	DeleteExpiredRefreshTokens() error
}
