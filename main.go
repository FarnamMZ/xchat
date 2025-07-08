package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"xchat/config"
	"xchat/handlers"
	"xchat/repository"
	"xchat/services"

	_ "github.com/lib/pq" // PostgresSQL driver
)

func main() {
	// Get config struct
	cfg, err := config.GetConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Couldn't get config from configfile: %s", err)
	}

	// Connect to the database.
	db, err := connectToDB(cfg)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %s", err)
	}

	// Ensure the database connection is alive.
	err = db.Ping()
	if err != nil {
		log.Fatalf("Couldn't ping database: %s", err)
	}

	// Initialize the repositories.
	usersRepository := repository.NewUsersRepository(db)
	jwtRepository := repository.NewJwtRepository(db)

	// Initialize the services.
	authService := services.NewAuthService(cfg, []byte(cfg.Jwt.Secret), usersRepository, jwtRepository)

	// Initialize the HTTP multiplexer with services.
	mux := handlers.NewMux(cfg, authService)

	// Set up the HTTP server.
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux.ServeMux,
	}

	// Start the server.
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Couldn't run the server: %s", err)
	}
}

// connectToDB establishes a connection to the PostgresSQL database and returns the database handle.
func connectToDB(cfg *config.Config) (*sql.DB, error) {
	dbCfg := cfg.Database
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbCfg.User, dbCfg.Password, dbCfg.Dbname))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
