package handlers

import (
	"net/http"
	"xchat/config"
)

// Mux is a custom HTTP multiplexer that handles routing for the application.
type Mux struct {
	ServeMux *http.ServeMux
	cfg      *config.Config
	as       AuthService
}

func NewMux(cfg *config.Config, as AuthService) *Mux {
	m := &Mux{
		ServeMux: http.NewServeMux(),
		cfg:      cfg,
		as:       as,
	}

	// Register the authentication handlers.
	m.ServeMux.HandleFunc("/signup", m.signup)
	m.ServeMux.HandleFunc("/login", m.login)
	m.ServeMux.HandleFunc("/logout", m.logout)

	return m
}
