package handlers

import (
	"errors"
	"net/http"
	"xchat/dto"
)

func (m *Mux) login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginReq
	err := m.decodeJSONBody(r, &req)
	if err != nil {
		m.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := m.as.Login(&req)
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
