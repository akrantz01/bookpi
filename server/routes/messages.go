package routes

import (
	"encoding/json"
	"github.com/akrantz01/bookpi/server/models"
	"github.com/akrantz01/bookpi/server/responses"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
)

// Routes for message management
func Messages(db *bolt.DB, router *mux.Router) {
	subrouter := router.PathPrefix("/chats/{chat}/messages").Subrouter()

	subrouter.HandleFunc("", messages(db))
}

// Operate on messages in a chat
func messages(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve chat id
		vars := mux.Vars(r)
		chat, ok := vars["chat"]
		if !ok {
			responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
			return
		}

		// Decode chat id
		id, err := uuid.FromString(chat)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "invalid chat id format")
			return
		}

		switch r.Method {
		case http.MethodGet:
			listMessages(w, r, id, db)

		case http.MethodPost:
			createMessage(w, r, id, db)

		default:
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// List the messages in a chat
func listMessages(w http.ResponseWriter, r *http.Request, chatId uuid.UUID, db *bolt.DB) {
	// Ensure chat exists
	chat, err := models.FindChat(chatId, db)
	if err != nil {
		log.Printf("ERROR: failed to query database for specified chat")
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	} else if chat == nil {
		responses.Error(w, http.StatusNotFound, "specified chat does not exist")
		return
	}

	// Get user
	self, err := models.FindUser(r.Header.Get("X-BPI-Username"), db)
	if err != nil {
		log.Printf("ERROR: failed to query database for requesting user: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	}

	// Check that user in chat
	for _, c := range self.Chats {
		if c == chatId.String() {
			responses.SuccessWithData(w, chat.Messages)
			return
		}
	}

	responses.Error(w, http.StatusForbidden, "user not in specified chat")
}

// Send a message to the specified chat
func createMessage(w http.ResponseWriter, r *http.Request, chatId uuid.UUID, db *bolt.DB) {
	// Validate initial request on headers and request body existence
	if r.Header.Get("Content-Type") != "application/json" {
		responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		responses.Error(w, http.StatusBadRequest, "request body must be present")
		return
	}

	// Parse and validate body fields
	var body struct{
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid json format for request body")
		return
	} else if body.Message == "" {
		responses.Error(w, http.StatusBadRequest, "field 'message' are required")
		return
	}

	// Ensure chat exists
	chat, err := models.FindChat(chatId, db)
	if err != nil {
		log.Printf("ERROR: failed to query database for specified chat")
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	} else if chat == nil {
		responses.Error(w, http.StatusNotFound, "specified chat does not exist")
		return
	}

	// Get user
	self, err := models.FindUser(r.Header.Get("X-BPI-Username"), db)
	if err != nil {
		log.Printf("ERROR: failed to query database for requesting user: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	}

	// Check that user in chat
	valid := false
	for _, c := range self.Chats {
		if c == chatId.String() {
			valid = true
			break
		}
	}
	if !valid {
		responses.Error(w, http.StatusForbidden, "user not in specified chat")
		return
	}

	// Add message chat
	chat.AddMessage(body.Message, self.Username)
	if err := chat.Save(db); err != nil {
		log.Printf("ERORR: failed to write new message to database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}

	responses.Success(w)
}
