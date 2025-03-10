//go:generate fileb0x b0x.yml

package main

import (
	"context"
	"errors"
	"github.com/akrantz01/bookpi/server/assets"
	"github.com/akrantz01/bookpi/server/models"
	"github.com/akrantz01/bookpi/server/responses"
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
		if !os.IsNotExist(err) {
			log.Fatalf("Failed to load environment variables: %v\n", err)
		}
	}
	cfg := loadEnv()

	// Delete old if resetting
	if cfg.Reset {
		if err := os.RemoveAll(cfg.FilesDirectory); err != nil {
			log.Fatalf("Failed to delete files directory: %v\n", err)
		}
		if err := os.Remove(cfg.Database); err != nil {
			log.Fatalf("Failed to delete database: %v\n", err)
		}
	}

	// Initialize file storage directory
	if err := os.MkdirAll(cfg.FilesDirectory, os.ModeDir|0755); err != nil {
		log.Fatalf("Failed to create files directory: %v\n", err)
	}

	// Initialize database
	db, err := bolt.Open(cfg.Database, 0600, &bolt.Options{Timeout: 5 * time.Second})
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
		for _, b := range [][]byte{models.BucketUsers, models.BucketSessions, models.BucketChats, models.BucketShares} {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return err
			}
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
	routes.Authentication(cfg.FilesDirectory, db, api)
	routes.Users(cfg.FilesDirectory, db, api)
	routes.Chats(db, api)
	routes.Messages(db, api)
	routes.Files(cfg.FilesDirectory, api)
	routes.Shares(cfg.FilesDirectory, db, api)

	// Register session middleware
	api.Use(sessionMiddleware(db))

	// Handle API errors
	api.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responses.Error(w, http.StatusNotFound, "route not found")
	})
	api.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	// Serve embedded files
	router.PathPrefix("/").Handler(assets.StaticServer)

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
			if err == errors.New("http: server closed") {
				return
			}

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
