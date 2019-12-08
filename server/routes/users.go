package routes

import (
	"encoding/base64"
	"encoding/json"
	"github.com/akrantz01/bookpi/server/hash"
	"github.com/akrantz01/bookpi/server/models"
	"github.com/akrantz01/bookpi/server/responses"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"
	"time"
)

// Routes for user management
func Users(filesDirectory string, db *bolt.DB, router *mux.Router) {
	subrouter := router.PathPrefix("/user").Subrouter()

	subrouter.HandleFunc("", selfUser(filesDirectory, db))
	subrouter.HandleFunc("/{username}", readUser("", db))
}

// Operate on the user in the session
func selfUser(filesDirectory string, db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve user from session
		id, _ := base64.URLEncoding.DecodeString(r.Header.Get("X-BPI-Session-Id"))
		session, err := models.FindSession(id, db)
		if err != nil {
			log.Printf("ERROR: failed to query database for user session: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to query database")
			return
		}

		switch r.Method {
		case http.MethodGet:
			readUser(session.User.Username, db)(w, r)

		case http.MethodPut:
			updateUser(w, r, session, db)

		case http.MethodDelete:
			deleteUser(w, r, session, filesDirectory, db)

		default:
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// Get a description of a user by their username
func readUser(user string, db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate initial request on method and path parameters
		vars := mux.Vars(r)
		if r.Method != http.MethodGet {
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		} else if _, ok := vars["username"]; !ok && user == "" {
			responses.Error(w, http.StatusBadRequest, "path parameter 'username' must be present")
			return
		} else if user != "" {
			vars["username"] = user
		}

		user, err := models.FindUser(vars["username"], db)
		if err != nil {
			log.Printf("ERROR: failed to query database for user: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to query database")
			return
		} else if user == nil {
			responses.Error(w, http.StatusNotFound, "specified user does not exist")
			return
		}

		responses.SuccessWithData(w, map[string]string{"name": user.Name, "username": user.Username})
	}
}

// Update a user's name or password
func updateUser(w http.ResponseWriter, r *http.Request, session *models.Session, db *bolt.DB) {
	if r.Header.Get("Content-Type") != "application/json" {
		responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		responses.Error(w, http.StatusBadRequest, "request body must be present")
		return
	}

	// Parse and validate body fields
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid json format for request body")
		return
	}

	// Update name if present
	if body.Name != "" {
		session.User.Name = body.Name
	}

	// Update password if present
	if body.Password != "" {
		// Validate password on requirements
		if len(body.Password) < 8 {
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

		h, err := hash.DefaultHash(body.Password)
		if err != nil {
			log.Printf("ERROR: failed to hash user password: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to hash password")
			return
		}

		session.User.Password = h
	}

	// Save user to database
	if err := session.User.Save(db); err != nil {
		log.Printf("ERROR: failed to write user updates to database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}

	responses.Success(w)
}

// Delete a user and invalidate their sessions
func deleteUser(w http.ResponseWriter, r *http.Request, _ *models.Session, filesDirectory string, db *bolt.DB) {
	// Get user from database
	user, err := models.FindUser(r.Header.Get("X-BPI-Username"), db)
	if err != nil {
		log.Printf("ERROR: failed to query user from database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	}

	// Batch delete sessions
	if err := db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("sessions"))

		// Delete all user's session ids
		for _, stringSID := range user.Sessions {
			sid, _ := base64.URLEncoding.DecodeString(stringSID)
			if err := bucket.Delete(sid); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Printf("ERROR: failed to delete sessions for user from database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to delete from database")
		return
	}

	// Delete the user's files
	if err := os.RemoveAll(filesDirectory + "/" + user.Username); err != nil {
		log.Printf("ERROR: failed to delete user file storage directory: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to delete directory")
		return
	}

	// Delete the user
	if err := user.Delete(db); err != nil {
		log.Printf("ERROR: failed to delete user from database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to delete from database")
		return
	}

        // Set empty cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "bp-id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   false,
		HttpOnly: true,
	})

	responses.Success(w)
}
