package models

import (
	"encoding/json"
	"github.com/akrantz01/bookpi/server/hash"
	bolt "go.etcd.io/bbolt"
)

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Create a new user
func NewUser(name, username, password string) (*User, error) {
	h, err := hash.DefaultHash(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:     name,
		Username: username,
		Password: h,
	}, nil
}

// Find a user by username
func FindUser(username string, db *bolt.DB) (*User, error) {
	var user User
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		// Decode user
		buf := bucket.Get([]byte(username))
		return json.Unmarshal(buf, &user)
	})

	switch err.(type) {
	case *json.SyntaxError:
		return nil, nil
	case nil:
		return &user, err
	default:
		return nil, err
	}
}

// Save a user to the database
func (u *User) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		// Marshal user data into bytes
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(u.Username), buf)
	})
}

// Check if a user's credentials are valid
func (u *User) Authenticate(username, password string) (bool, error) {
	if u.Username != username {
		return false, nil
	}

	return hash.Verify(password, u.Password)
}

// Delete the user from the database
func (u *User) Delete(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		return bucket.Delete([]byte(u.Username))
	})
}
