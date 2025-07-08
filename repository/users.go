package repository

import (
	"database/sql"
	"errors"
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

func (ur *usersRepository) SaveUser(user *models.User) error {
	query := `INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4) returning id`
	err := ur.db.QueryRow(query, user.Username, user.Password, user.Email, user.Role).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *usersRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, email, role FROM users WHERE username = $1`
	row := ur.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}
	return &user, nil
}

func (ur *usersRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, username, password, email, role FROM users WHERE id = $1`
	row := ur.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}
	return &user, nil
}
