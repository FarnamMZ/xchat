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

	at, rt, err := m.as.Signup(&req)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			m.respondWithError(w, httpErr.StatusCode, httpErr.Message)
			return
		}
		m.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

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
}
