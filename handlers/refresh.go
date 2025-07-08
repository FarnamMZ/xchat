package handlers

import (
	"errors"
	"net/http"
	"time"
)

func (m *Mux) refreshToken(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie
	cookie, err := r.Cookie("refresh")
	if err != nil {
		m.respondWithError(w, http.StatusUnauthorized, "Refresh token not found")
		return
	}

	refreshToken := cookie.Value
	if refreshToken == "" {
		m.respondWithError(w, http.StatusUnauthorized, "Refresh token is empty")
		return
	}

	// Use the auth service to refresh the token
	at, rt, err := m.as.RefreshToken(refreshToken)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			m.respondWithError(w, httpErr.StatusCode, httpErr.Message)
			return
		}
		m.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Set new tokens as cookies
	atAge := m.cfg.Jwt.TokenAge
	atCookie := http.Cookie{
		Name:     "token",
		Value:    at,
		MaxAge:   int(time.Minute) * atAge,
		Expires:  time.Now().Add(time.Minute * time.Duration(atAge)),
		Secure:   true,
		HttpOnly: true,
	}

	rtAge := m.cfg.Jwt.RefreshAge
	rtCookie := http.Cookie{
		Name:     "refresh",
		Value:    rt,
		MaxAge:   int(time.Minute) * rtAge,
		Expires:  time.Now().Add(time.Minute * time.Duration(rtAge)),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &rtCookie)

	w.WriteHeader(http.StatusOK)
}
