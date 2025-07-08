package models

import "time"

type RefreshToken struct {
	JTI       string    `json:"jti"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
