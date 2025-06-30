package handlers

import "net/http"

func (m *Mux) login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login\n"))
}
