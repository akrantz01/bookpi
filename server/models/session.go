package models

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
)

type Session struct {
	Id   uuid.UUID `json:"-"`
	User User      `json:"user"`
}

// Create a new session
func NewSession(user User) Session {
	return Session{
		Id:   uuid.NewV4(),
		User: user,
	}
}

// Retrieve session data from database
func FindSession(id uuid.UUID, db *bolt.DB) (*Session, error) {
	var session Session
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("sessions"))
		raw := bucket.Get(id.Bytes())
		return json.Unmarshal(raw, &session)
	})
	session.Id = id
	return &session, err
}

// Save the session to the database
func (s *Session) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("sessions"))

		// Marshal into JSON
		buf, err := json.Marshal(s)
		if err != nil {
			return err
		}

		return bucket.Put(s.Id.Bytes(), buf)
	})
}

// Delete the session from the database
func (s *Session) Delete(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("sessions"))
		return bucket.Delete(s.Id.Bytes())
	})
}
