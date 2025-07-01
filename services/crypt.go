package services

import (
	"beta/handlers"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type cryptService struct{}

// NewCryptService creates a new instance of CryptService
func NewCryptService() handlers.CryptService {
	return &cryptService{}
}

func (cs *cryptService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (cs *cryptService) ComparePassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func (cs *cryptService) IsPasswordSecure(password string) bool {
	// password should be at least 8 characters long and 50 characters max
	// password should contain at least one uppercase letter, one lowercase letter, one number, and one special character
	if len(password) < 8 || len(password) > 50 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	specialChars := "!@#$%^&*()-_=+[]{}|;:',.<>?/~`"

	for _, letter := range password {
		if letter >= 'A' && letter <= 'Z' {
			hasUpper = true
		} else if letter >= 'a' && letter <= 'z' {
			hasLower = true
		} else if letter >= '0' && letter <= '9' {
			hasNumber = true
		} else if strings.Contains(specialChars, string(letter)) {
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
