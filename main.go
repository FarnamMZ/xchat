package main

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgresSQL driver
	"log"
	"net/http"
	"xchat/handlers"
	"xchat/repository"
	"xchat/services"
)

func main() {
	// Connect to the database.
	db, err := connectToDB()
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

	// Initialize the services.
	authService := services.NewAuthService([]byte("my_secret_key"), usersRepository)
	cryptService := services.NewCryptService()

	// Initialize the HTTP multiplexer with services.
	mux := handlers.NewMux(authService, cryptService)

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

// connectToDB establishes a connection to the PostgreSQL database and returns the database handle.
func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=farnammrz password=1122 dbname=beta sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
