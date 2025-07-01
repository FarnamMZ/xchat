package repository

import (
	"database/sql"
	"xchat/models"
	"xchat/services"
)

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) services.UsersRepository {
	return &usersRepository{db: db}
}

func (ur *usersRepository) UserExistsByUsername(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := ur.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (ur *usersRepository) UserExistsByEmail(emil string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := ur.db.QueryRow(query, emil).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (ur *usersRepository) SaveUser(user models.User) error {
	query := `INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4)`
	_, err := ur.db.Exec(query, user.Username, user.Password, user.Email, user.Role)
	if err != nil {
		return err
	}
	return nil
}
