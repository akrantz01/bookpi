package main

import (
	"context"
	"github.com/akrantz01/bookpi/server/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	db, err := bolt.Open(cfg.Database, 0600, &bolt.Options{Timeout: 5*time.Second})
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close database: %v\n", err)
		}
	}()

	// Create database buckets if not exist
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("users")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("sessions")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	// Listen for OS signals
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Add server router
	router := mux.NewRouter()

	// Register API routes
	api := router.PathPrefix("/api").Subrouter()
	routes.Authentication(db, api)
	routes.Users(db, api)

	// Register session middleware
	router.Use(sessionMiddleware(db))

	// Setup server
	server := http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      applyWrappers(router),
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
