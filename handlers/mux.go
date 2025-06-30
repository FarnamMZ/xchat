package handlers

import (
	"net/http"
)

// Mux is a custom HTTP multiplexer that handles routing for the application.
type Mux struct {
	ServeMux *http.ServeMux
	as       AuthService
}

func NewMux(as AuthService) *Mux {
	m := &Mux{
		ServeMux: http.NewServeMux(),
		as:       as,
	}

	// Register the authentication handlers.
	m.ServeMux.HandleFunc("/signup", m.signup)
	m.ServeMux.HandleFunc("/login", m.login)
	m.ServeMux.HandleFunc("/logout", m.logout)

	return m
}
