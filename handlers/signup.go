package handlers

import (
	"errors"
	"net/http"
	"time"
	"xchat/dto"
)

func (m *Mux) signup(w http.ResponseWriter, r *http.Request) {
	var req dto.SignupReq
	err := m.decodeJSONBody(r, &req)
	if err != nil {
		m.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := m.as.Signup(&req)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			m.respondWithError(w, httpErr.StatusCode, httpErr.Message)
			return
		}
		m.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	age := m.cfg.Jwt.TokenAge
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   int(time.Minute) * age,
		Expires:  time.Now().Add(time.Minute * time.Duration(age)),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
