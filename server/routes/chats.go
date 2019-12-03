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

// Routes for chat management
func Chats(db *bolt.DB, router *mux.Router) {
	subrouter := router.PathPrefix("/chats").Subrouter()

	subrouter.HandleFunc("", allChats(db))
	subrouter.HandleFunc("/{chat}", specificChat(db))
}

// Operate on all a user's chats
func allChats(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listChats(w, r, db)

		case http.MethodPost:
			createChat(w, r, db)

		default:
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// Operate on a specific chat
func specificChat(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
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
			readChat(w, r, id, db)

		case http.MethodDelete:
			deleteChat(w, r, id, db)

		default:
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// Get a list of the user's chats
func listChats(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// Get self
	self, err := models.FindUser(r.Header.Get("X-BPI-Username"), db)
	if err != nil {
		log.Printf("ERROR: failed to query database for requesting user")
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	}

	// Return empty array if none
	if len(self.Chats) == 0 {
		responses.SuccessWithData(w, []string{})
		return
	}

	responses.SuccessWithData(w, self.Chats)
}

// Create a chat between two users
func createChat(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// Validate initial request headers, and body existence
	if r.Header.Get("Content-Type") != "application/json" {
		responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		responses.Error(w, http.StatusBadRequest, "request body must be present")
		return
	}

	// Parse and validate body fields
	var body struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid json format for request body")
		return
	} else if body.To == "" || body.Message == "" {
		responses.Error(w, http.StatusBadRequest, "fields 'name' and 'message' are required")
		return
	}

	// Get self
	self, err := models.FindUser(r.Header.Get("X-BPI-Username"), db)
	if err != nil {
		log.Printf("ERROR: failed to query database for self: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	}

	// Check if specified user exists
	recipient, err := models.FindUser(body.To, db)
	if err != nil {
		log.Printf("ERROR: failed to query database for recipient user existence: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to query database")
		return
	} else if recipient == nil {
		responses.Error(w, http.StatusBadRequest, "specified recipient does not exist")
		return
	}

	// Create the chat
	chat := models.NewChat(self.Username, recipient.Username, body.Message)
	if err := chat.Save(db); err != nil {
		log.Printf("ERROR: failed to write chat to database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}

	// Add chat id to each user
	self.AddChat(chat.Id.String())
	recipient.AddChat(chat.Id.String())

	// Save users
	if err := self.Save(db); err != nil {
		log.Printf("ERROR: failed to write chat id to requesting user: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}
	if err := recipient.Save(db); err != nil {
		log.Printf("ERROR: failed to write chat id to recipient user: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}

	responses.Success(w)
}

// Describe a user's chat
func readChat(w http.ResponseWriter, r *http.Request, chatId uuid.UUID, db *bolt.DB) {
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
			responses.SuccessWithData(w, map[string]string{
				"user1":        chat.User1,
				"user2":        chat.User2,
				"last_message": chat.Messages[len(chat.Messages)-1],
			})
			return
		}
	}

	responses.Error(w, http.StatusForbidden, "user not in specified chat")
}

// Delete a user's chat
func deleteChat(w http.ResponseWriter, r *http.Request, chatId uuid.UUID, db *bolt.DB) {
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

	// Get second user
	var user *models.User
	if self.Username == chat.User1 {
		user, err = models.FindUser(chat.User2, db)
		if err != nil {
			log.Printf("ERROR: failed to query database for second user in chat: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to query database")
			return
		}
	} else {
		user, err = models.FindUser(chat.User1, db)
		if err != nil {
			log.Printf("ERROR: failed to query database for second user in chat: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to query database")
			return
		}
	}

	// Remove chat from second user if exists
	if user != nil {
		user.RemoveChat(chat.Id.String())
		if err := user.Save(db); err != nil {
			log.Printf("ERROR: failed to write user information to database: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to write to database")
			return
		}
	}

	// Remove chat for requesting user
	self.RemoveChat(chat.Id.String())
	if err := self.Save(db); err != nil {
		log.Printf("ERROR: failed to write user information to database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to write to database")
		return
	}

	// Remove chat entirely
	if err := chat.Delete(db); err != nil {
		log.Printf("ERROR: failed to delete chat from database: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to delete from database")
		return
	}

	responses.Success(w)
}
