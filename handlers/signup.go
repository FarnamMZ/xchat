package handlers

import (
	"errors"
	"net/http"
	"xchat/dto"
)

func (m *Mux) signup(w http.ResponseWriter, r *http.Request) {
	var req dto.SignupReq
	err := m.decodeJSONBody(r, &req)
	if err != nil {
		m.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// check if password is strong
	if !m.cs.IsPasswordSecure(req.Password) {
		httpErr := ErrInsecurePassword
		m.respondWithError(w, httpErr.StatusCode, httpErr.Message)
		return
	}

	// hash password
	req.Password, err = m.cs.HashPassword(req.Password)
	if err != nil {
		m.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := m.as.Signup(req)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			m.respondWithError(w, httpErr.StatusCode, httpErr.Message)
			return
		}
		m.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	m.respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
