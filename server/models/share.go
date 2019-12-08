package models

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
)

type Share struct {
	Path string   `json:"-"`
	To   []string `json:"to"`
}

// Create a new file share
func NewShare(file string) *Share {
	return &Share{
		Path: file,
		To:   []string{},
	}
}

// Find a file share by id
func FindShare(path string, db *bolt.DB) (*Share, error) {
	var share Share
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketShares)

		// Decode share
		buf := bucket.Get([]byte(path))
		return json.Unmarshal(buf, &share)
	})

	switch err.(type) {
	case *json.SyntaxError:
		return nil, nil
	case nil:
		share.Path = path
		return &share, nil
	default:
		return nil, err
	}
}

// Save a share to the database
func (s *Share) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketShares)

		// Marshal share data into bytes
		buf, err := json.Marshal(s)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(s.Path), buf)
	})
}

// Add user to shared file
func (s *Share) AddUser(user string) {
	s.To = append(s.To, user)
}

// Remove user from shared file
func (s *Share) RemoveUser(user string) {
	for i, u := range s.To {
		if u == user {
			s.To = append(s.To[:i], s.To[i+1:]...)
			break
		}
	}
}

// Delete a share from the database
func (s *Share) Delete(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketShares)
		return bucket.Delete([]byte(s.Path))
	})
}
