package handlers

import (
	"net/http"
	"time"
)

func (m *Mux) logout(w http.ResponseWriter, r *http.Request) {
	// clear the session cookie
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
