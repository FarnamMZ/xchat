package services

import "xchat/models"

type UsersRepository interface {
	UserExistsByUsername(username string) (bool, error)
	UserExistsByEmail(email string) (bool, error)
	SaveUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
}
