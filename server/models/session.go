package models

import (
	"crypto/rand"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"io"
)

type Session struct {
	Id   []byte `json:"-"`
	User User   `json:"user"`
}

// Create a new session
func NewSession(user User) Session {
	b := make([]byte, 64)
	_, _ = io.ReadFull(rand.Reader, b)

	return Session{
		Id:   b,
		User: user,
	}
}

// Retrieve session data from database
func FindSession(id []byte, db *bolt.DB) (*Session, error) {
	var session Session
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("sessions"))
		raw := bucket.Get(id)
		return json.Unmarshal(raw, &session)
	})

	switch err.(type) {
	case *json.SyntaxError:
		return nil, nil
	case nil:
		session.Id = id
		return &session, nil
	default:
		return nil, err
	}
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

		return bucket.Put(s.Id, buf)
	})
}

// Delete the session from the database
func (s *Session) Delete(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("sessions"))
		return bucket.Delete(s.Id)
	})
}
