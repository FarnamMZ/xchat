package handlers

import (
	"net/http"
	"time"

	"xchat/models"

	"github.com/golang-jwt/jwt"
)

func (m *Mux) logout(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie to invalidate it
	refreshCookie, err := r.Cookie("refresh")
	if err == nil && refreshCookie.Value != "" {
		// Parse the refresh token to get the JTI
		token, err := jwt.ParseWithClaims(refreshCookie.Value, &models.RTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.Jwt.Secret), nil
		})
		if err == nil {
			if claims, ok := token.Claims.(*models.RTClaims); ok && token.Valid {
				// Delete the refresh token from database
				m.as.DeleteRefreshToken(claims.Id)
			}
		}
	}

	// Clear the access token cookie
	atCookie := http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		Secure:   true,
		HttpOnly: true,
	}

	// Clear the refresh token cookie
	rtCookie := http.Cookie{
		Name:     "refresh",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &rtCookie)
}
