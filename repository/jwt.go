package repository

import (
	"database/sql"
	"time"
	"xchat/models"
)

type jwtRepository struct {
	db *sql.DB
}

func NewJwtRepository(db *sql.DB) *jwtRepository {
	return &jwtRepository{db: db}
}

func (jr *jwtRepository) SaveRefreshToken(token *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (jti, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4)`
	_, err := jr.db.Exec(query, token.JTI, token.UserID, token.ExpiresAt, token.CreatedAt)
	return err
}

func (jr *jwtRepository) GetRefreshToken(jti string) (*models.RefreshToken, error) {
	query := `SELECT jti, user_id, expires_at, created_at FROM refresh_tokens WHERE jti = $1`
	row := jr.db.QueryRow(query, jti)

	var token models.RefreshToken
	err := row.Scan(&token.JTI, &token.UserID, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func (jr *jwtRepository) DeleteRefreshToken(jti string) error {
	query := `DELETE FROM refresh_tokens WHERE jti = $1`
	_, err := jr.db.Exec(query, jti)
	return err
}

func (jr *jwtRepository) DeleteExpiredRefreshTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	_, err := jr.db.Exec(query, time.Now())
	return err
}
