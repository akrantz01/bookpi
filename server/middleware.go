package main

import (
	"encoding/base64"
	"github.com/akrantz01/bookpi/server/models"
	"github.com/akrantz01/bookpi/server/responses"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"
)

// Apply the wrapper functions
func applyWrappers(router *mux.Router) http.Handler {
	logging := handlers.CombinedLoggingHandler(os.Stdout, router)
	corsEnabled := cors.AllowAll().Handler(logging)
	return corsEnabled
}

func sessionMiddleware(db *bolt.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow if authenticating or registering
			if r.RequestURI == "/api/auth/login" || r.RequestURI == "/api/auth/register" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if cookie exists
			cookie, err := r.Cookie("bp-id")
			if err != nil {
				responses.Error(w, http.StatusUnauthorized, "no session present")
				return
			}

			// Parse id from cookie
			id, err := base64.URLEncoding.DecodeString(cookie.Value)
			if err != nil {
				responses.Error(w, http.StatusUnauthorized, "invalid session id format")
				return
			}

			// Retrieve cookie from database
			session, err := models.FindSession(id, db)
			if err != nil {
				log.Printf("ERROR: failed to query database for session: %v\n", err)
				responses.Error(w, http.StatusInternalServerError, "failed to query database")
				return
			} else if session == nil {
				responses.Error(w, http.StatusUnauthorized, "invalid session id")
				return
			}

			// Set data from session to headers
			r.Header.Set("X-BPI-Session-Id", cookie.Value)
			r.Header.Set("X-BPI-Username", session.User.Username)
			r.Header.Set("X-BPI-Name", session.User.Name)

			next.ServeHTTP(w, r)
		})
	}
}
