package models

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
)

type Chat struct {
	Id uuid.UUID `json:"-"`
	User1 string `json:"user1"`
	User2 string `json:"user2"`
	Messages []string `json:"messages"`
}

// Create a new chat
func NewChat(user1, user2, initialMessage string) *Chat {
	return &Chat{
		Id:       uuid.NewV4(),
		User1:    user1,
		User2:    user2,
		Messages: []string{initialMessage},
	}
}

// Find a chat by id
func FindChat(id uuid.UUID, db *bolt.DB) (*Chat, error) {
	var chat Chat
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("chats"))

		// Decode chat
		buf := bucket.Get(id.Bytes())
		return json.Unmarshal(buf, &chat)
	})

	switch err.(type) {
	case *json.SyntaxError:
		return nil, nil
	case nil:
		return &chat, nil
	default:
		return nil, err
	}
}

// Save a chat to the database
func (c *Chat) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("chats"))

		// Marshal chat data into bytes
		buf, err := json.Marshal(c)
		if err != nil {
			return err
		}

		return bucket.Put(c.Id.Bytes(), buf)
	})
}

// Add a message to the chat
func (c *Chat) AddMessage(message, from string) {
	c.Messages = append(c.Messages, from+":"+message)
}

// Remove a message from the chat
func (c *Chat) RemoveMessage(index int) {
	c.Messages = append(c.Messages[:index], c.Messages[index+1:]...)
}

// Delete a chat from the database
func (c *Chat) Delete(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("chats"))
		return bucket.Delete(c.Id.Bytes())
	})
}
