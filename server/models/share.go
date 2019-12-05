package models

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
)

type Share struct {
	Id uuid.UUID `json:"-"`
	File string `json:"file"`
	From string `json:"from"`
}

// Create a new file share
func NewShare(file, user string) *Share {
	return &Share{
		Id: uuid.NewV4(),
		From: user,
		File: file,
	}
}

// Find a file share by id
func FindShare(id uuid.UUID, db *bolt.DB) (*Share, error) {
	var share Share
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketShares)

		// Decode share
		buf := bucket.Get(id.Bytes())
		return json.Unmarshal(buf, &share)
	})

	switch err.(type) {
	case *json.SyntaxError:
		return nil, nil
	case nil:
		share.Id = id
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

		return bucket.Put(s.Id.Bytes(), buf)
	})
}

// Delete a share from the database
func (s *Share) Delete(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketShares)
		return bucket.Delete(s.Id.Bytes())
	})
}
