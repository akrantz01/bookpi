package models

import (
	"crypto/rand"
	"encoding/base64"
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
	// Generate id
	b := make([]byte, 64)
	_, _ = io.ReadFull(rand.Reader, b)

	// Add session id to user
	user.Sessions = append(user.Sessions, base64.URLEncoding.EncodeToString(b))

	return Session{
		Id:   b,
		User: user,
	}
}

// Retrieve session data from database
func FindSession(id []byte, db *bolt.DB) (*Session, error) {
	var session Session
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketSessions)
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
	// Save the session itself
	if err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketSessions)

		// Marshal into JSON
		buf, err := json.Marshal(s)
		if err != nil {
			return err
		}

		return bucket.Put(s.Id, buf)
	}); err != nil {
		return err
	}

	// Update the user
	return s.User.Save(db)
}

// Delete the session from the database
func (s *Session) Delete(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketSessions)
		return bucket.Delete(s.Id)
	})
}
