package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize configuration from environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load environment variables: %v\n", err)
	}
	cfg := loadEnv()

	// Initialize database
	db, err := bolt.Open(cfg.Database, 0600, nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close database: %v\n", err)
		}
	}()

	// Listen for OS signals
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Add server router
	router := mux.NewRouter()

	// Setup server
	server := http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      cors.AllowAll().Handler(handlers.CombinedLoggingHandler(os.Stdout, router)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start in separate goroutine
	go func() {
		log.Printf("Listening on %s:%s...\n", cfg.Host, cfg.Port)

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Error while running server: %v\n", err)
		}
	}()

	// Block for shutdown signal
	<-shutdown

	// Create the shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown th server: %v\n", err)
	}

	log.Println("server is shutdown, goodbye")
}
