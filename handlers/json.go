package handlers

import (
	"encoding/json"
	"net/http"
)

func (m *Mux) respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (m *Mux) respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func (m *Mux) decodeJSONBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return http.ErrBodyNotAllowed
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents unknown fields in the JSON body
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	return nil
}
