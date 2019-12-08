package routes

import (
	"encoding/base64"
	"encoding/json"
	"github.com/akrantz01/bookpi/server/models"
	"github.com/akrantz01/bookpi/server/responses"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

var (
	regexUsername  = regexp.MustCompile("^([a-zA-Z0-9]+)$")
	regexLowercase = regexp.MustCompile("[a-z]+")
	regexUppercase = regexp.MustCompile("[A-Z]")
	regexNumeric   = regexp.MustCompile("[0-9]+")
	regexSpecial   = regexp.MustCompile("[!-_]+")
)

// Handle user authentication
func Authentication(filesDirectory string, db *bolt.DB, router *mux.Router) {
	subrouter := router.PathPrefix("/auth").Subrouter()

	subrouter.HandleFunc("/register", register(filesDirectory, db))
	subrouter.HandleFunc("/login", login(db))
	subrouter.HandleFunc("/logout", logout(db))
}

// Handle user registration
func register(filesDirectory string, db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate initial request on method, headers, and body existence
		if r.Method != http.MethodPost {
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		} else if r.Header.Get("Content-Type") != "application/json" {
			responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
			return
		} else if r.Body == nil {
			responses.Error(w, http.StatusBadRequest, "request body must be present")
			return
		}

		// Parse and validate body fields
		var body struct {
			Name     string `json:"name"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			responses.Error(w, http.StatusBadRequest, "invalid json format for request body")
			return
		} else if body.Name == "" || body.Username == "" || body.Password == "" {
			responses.Error(w, http.StatusBadRequest, "fields 'name', 'username', and 'password' are required")
			return
		} else if len(body.Username) < 3 {
			responses.Error(w, http.StatusBadRequest, "field 'username' must be at least 3 characters")
			return
		} else if !regexUsername.MatchString(body.Username) {
			responses.Error(w, http.StatusBadRequest, "field 'username' must only contain lowercase characters")
			return
		} else if len(body.Password) < 8 {
			responses.Error(w, http.StatusBadRequest, "field 'password' must be at least 8 characters")
			return
		} else if !regexLowercase.MatchString(body.Password) {
			responses.Error(w, http.StatusBadRequest, "field 'password' must contain a lowercase character")
			return
		} else if !regexUppercase.MatchString(body.Password) {
			responses.Error(w, http.StatusBadRequest, "field 'password' must contain a uppercase character")
			return
		} else if !regexNumeric.MatchString(body.Password) {
			responses.Error(w, http.StatusBadRequest, "field 'password' must contain a numeric character")
			return
		} else if !regexSpecial.MatchString(body.Password) {
			responses.Error(w, http.StatusBadRequest, "field 'password' must contain a special character")
			return
		}

		// Check if user already exists
		u, err := models.FindUser(body.Username, db)
		if err != nil {
			log.Printf("ERROR: failed to query database for user: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to query database")
			return
		} else if u != nil {
			responses.Error(w, http.StatusConflict, "specified username is already in use")
			return
		}

		// Create the user
		u, err = models.NewUser(body.Name, body.Username, body.Password)
		if err != nil {
			log.Printf("ERROR: failed to hash user password: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to hash password")
			return
		}

		// Create user file directory
		if err := os.Mkdir(filesDirectory+"/"+u.Username, os.ModeDir|0755); err != nil {
			log.Printf("ERROR: failed to create user directory for file storage: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to create directory")
			return
		}

		// Save to database
		if err := u.Save(db); err != nil {
			log.Printf("ERROR: failed to write user information to database")
			responses.Error(w, http.StatusInternalServerError, "failed to write to database")
			return
		}

		responses.Success(w)
	}
}

// Handle user login
func login(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate initial request on method, headers, and body existence
		if r.Method != http.MethodPost {
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		} else if r.Header.Get("Content-Type") != "application/json" {
			responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
			return
		} else if r.Body == nil {
			responses.Error(w, http.StatusBadRequest, "request body must be present")
			return
		}

		// Parse and validate body fields
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			responses.Error(w, http.StatusBadRequest, "invalid json format for request body")
			return
		} else if body.Username == "" || body.Password == "" {
			responses.Error(w, http.StatusBadRequest, "fields 'username' and 'password' are required")
			return
		}

		// Find user
		user, err := models.FindUser(body.Username, db)
		if err != nil {
			log.Printf("ERROR: failed to query database for user: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to query database")
			return
		} else if user == nil {
			responses.Error(w, http.StatusUnauthorized, "invalid username or password")
			return
		}

		// Authenticate user
		if valid, err := user.Authenticate(body.Username, body.Password); err != nil {
			log.Printf("ERROR: failed to verify password against hash: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to verify password")
			return
		} else if !valid {
			responses.Error(w, http.StatusUnauthorized, "invalid username or password")
			return
		}

		// Create new session
		session := models.NewSession(*user)
		if err := session.Save(db); err != nil {
			log.Printf("ERROR: failed to save session to database: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to write to database")
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "bp-id",
			Value:    base64.URLEncoding.EncodeToString(session.Id),
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   false,
			HttpOnly: true,
		})

		responses.Success(w)
	}
}

// Handle user logout
func logout(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate initial request on method, headers, and body existence
		if r.Method != http.MethodGet {
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// Get session from database
		id, _ := base64.URLEncoding.DecodeString(r.Header.Get("X-BPI-Session-Id"))
		session, err := models.FindSession(id, db)
		if err != nil {
			log.Printf("ERROR: failed to query database for session: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to query database")
			return
		}

		// Delete session
		if err := session.Delete(db); err != nil {
			log.Printf("ERROR: failed to delete session from database: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to delete from database")
			return
		}

                // Set empty cookie
                http.SetCookie(w, &http.Cookie{
                        Name: "bp-id",
                        Value: "",
                        Path: "/",
                        Expires: time.Unix(0, 0),
                        Secure: false,
                        HttpOnly: true,
                })

		responses.Success(w)
	}
}
