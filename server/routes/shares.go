package routes

import (
	"encoding/json"
	"github.com/akrantz01/bookpi/server/models"
	"github.com/akrantz01/bookpi/server/responses"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Shares(filesDirectory string, db *bolt.DB, router *mux.Router) {
	subrouter := router.PathPrefix("/shares").Subrouter()

	subrouter.HandleFunc("", allShares(filesDirectory, db))
	subrouter.PathPrefix("/{user}/").HandlerFunc(specificShare(filesDirectory, db))
}

// Operate on all a user's shares
func allShares(filesDirectory string, db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listShares(w, r, db)

		case http.MethodPost:
			createShare(w, r, filesDirectory, db)

		default:
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// Operate on a specific user's share
func specificShare(filesDirectory string, db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			downloadShare(w, r, filesDirectory, db)

		case http.MethodDelete:
			deleteShare(w, r, db)

		default:
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// Create a link shared file
func createShare(w http.ResponseWriter, r *http.Request, filesDirectory string, db *bolt.DB) {
	// Validate initial request on headers and body existence
	if r.Header.Get("Content-Type") != "application/json" {
		responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		responses.Error(w, http.StatusBadRequest, "request body must be present")
		return
	}

	// Parse and validate body fields
	var body struct{
		File string `json:"file"`
		To string `json:"to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid json format for request body")
		return
	} else if body.File == "" || body.To == "" {
		responses.Error(w, http.StatusBadRequest, "fields 'file' and 'to' must be present")
		return
	}

	// Create the path
	namespacedPath := filepath.Join(r.Header.Get("X-BPI-Username"), strings.Replace(body.File, "../", "", -1))
	path := filepath.Join(filesDirectory, namespacedPath)

	// Ensure requested file exists
	if info, err := os.Stat(path); os.IsNotExist(err) {
		responses.Error(w, http.StatusNotFound, "specified file/directory does not exist")
		return
	} else if info.IsDir() {
		responses.Error(w, http.StatusBadRequest, "cannot share directory")
		return
	} else if err != nil {
		log.Printf("ERROR: failed to check file existence: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to stat file")
		return
	}

	// Ensure user to exists
	user, err := models.FindUser(body.To, db)
	if err != nil {
		log.Printf("ERROR: failed to query database for user existence: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	} else if user == nil {
		responses.Error(w, http.StatusNotFound, "specified user does not exist")
		return
	}

	// Check if share already exists
	share, err := models.FindShare(namespacedPath, db)
	if err != nil {
		log.Printf("ERROR: failed to query database for share existence: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	} else if share == nil {
		// Create a new share
		share = models.NewShare(namespacedPath)
	}

	// Ensure not already in share
	for _, user := range share.To {
		if user == body.To {
			responses.Error(w, http.StatusBadRequest, "already shared with user")
			return
		}
	}

	// Add user to share
	share.AddUser(body.To)

	// Save to database
	if err := share.Save(db); err != nil {
		log.Printf("ERROR: failed to write link share to database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}

	// Add to user
	user.AddShare(namespacedPath)
	if err := user.Save(db); err != nil {
		log.Printf("ERROR: failed to write updated user to database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}

	responses.Success(w)
}

// Get all a user's link shared files
func listShares(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// Get self
	self, err := models.FindUser(r.Header.Get("X-BPI-Username"), db)
	if err != nil {
		log.Printf("ERROR: failed to query database for requesting user: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	}

	// Send empty array instead of null
	if len(self.Shares) == 0 {
		responses.SuccessWithData(w, []string{})
		return
	}

	responses.SuccessWithData(w, self.Shares)
}

// Download a shared file
func downloadShare(w http.ResponseWriter, r *http.Request, filesDirectory string, db *bolt.DB) {
	// Validate initial request on method and path parameters
	vars := mux.Vars(r)
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if _, ok := vars["user"]; !ok {
		responses.Error(w, http.StatusBadRequest, "path parameter 'user' is required")
		return
	}

	// Assemble paths
	namespacedPath := filepath.Join(vars["user"], strings.TrimPrefix(r.URL.Path, "/api/shares/"+vars["user"]))
	path := filepath.Join(filesDirectory, namespacedPath)

	// Ensure path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		responses.Error(w, http.StatusNotFound, "specified file/directory does not exist")
		return
	} else if err != nil {
		log.Printf("ERROR: failed to check for file existence: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to stat file")
		return
	}

	// Ensure the share exists
	share, err := models.FindShare(namespacedPath, db)
	if err != nil {
		log.Printf("ERROR: failed to query database for share existence: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	} else if share == nil {
		responses.Error(w, http.StatusNotFound, "specified share does not exist")
		return
	}

	// Get file description if query param
	if r.URL.Query().Get("describe") != "" && vars["user"] == r.Header.Get("X-BPI-Username") {
		// Only allow sharer
		if vars["user"] != r.Header.Get("X-BPI-Username") {
			responses.Error(w, http.StatusForbidden, "cannot describe file")
			return
		}

		responses.SuccessWithData(w, share.To)
		return
	}

	// Ensure user in share
	valid := false
	for _, user := range share.To {
		if user == r.Header.Get("X-BPI-Username") {
			valid = true
			break
		}
	}
	if !valid {
		responses.Error(w, http.StatusForbidden, "file not shared with requesting user")
		return
	}

	http.ServeFile(w, r, path)
}

// Delete the entire share or a specific user from a share
func deleteShare(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["user"]; !ok {
		responses.Error(w, http.StatusBadRequest, "path parameter 'user' must be present")
		return
	}

	// Ensure user owns link share
	if r.Header.Get("X-BPI-Username") != vars["user"] {
		responses.Error(w, http.StatusForbidden, "requesting user does not own link share")
		return
	}

	// Assemble paths
	namespacedPath := filepath.Join(vars["user"], strings.TrimPrefix(r.URL.Path, "/api/shares/"+vars["user"]))

	// Ensure share exists
	share, err := models.FindShare(namespacedPath, db)
	if err != nil {
		log.Printf("ERORR: failed to query database for link share existence: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	} else if share == nil {
		responses.Error(w, http.StatusNotFound, "specified shared link does not exist")
		return
	}

	// Check if removing user or share
	if user := r.URL.Query().Get("user"); user != "" {
		share.RemoveUser(user)

		// Update database
		if err := share.Save(db); err != nil {
			log.Printf("ERROR: failed to write updates to database: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to write to database")
			return
		}

		responses.Success(w)
		return
	}

	// Delete entire share
	if err := share.Delete(db); err != nil {
		log.Printf("ERROR: failed to delete share from database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}

	responses.Success(w)
}
