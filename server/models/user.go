package models

import (
	"encoding/json"
	"github.com/akrantz01/bookpi/server/hash"
	bolt "go.etcd.io/bbolt"
)

type User struct {
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Sessions []string `json:"sessions"`
	Chats    []string `json:"chats"`
	Shares   []string `json:"shares"`
}

// Shares stores:
//   key: id
//   - path -> path to file
//   - user -> user shared by

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
		Sessions: []string{},
		Chats:    []string{},
		Shares:   []string{},
	}, nil
}

// Find a user by username
func FindUser(username string, db *bolt.DB) (*User, error) {
	var user User
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketUsers)

		// Decode user
		buf := bucket.Get([]byte(username))
		return json.Unmarshal(buf, &user)
	})

	switch err.(type) {
	case *json.SyntaxError:
		return nil, nil
	case nil:
		return &user, nil
	default:
		return nil, err
	}
}

// Save a user to the database
func (u *User) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketUsers)

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

// Associate a chat with the user
func (u *User) AddChat(id string) {
	u.Chats = append(u.Chats, id)
}

// Disassociate a user from the chat
func (u *User) RemoveChat(id string) {
	for i, chat := range u.Chats {
		if chat == id {
			u.Chats = append(u.Chats[:i], u.Chats[i+1:]...)
			break
		}
	}
}

// Associate a shared link with the user
func (u *User) AddShare(path string) {
	u.Shares = append(u.Shares, path)
}

// Disassociate a user from a shared link
func (u *User) RemoveShare(path string) {
	for i, share := range u.Shares {
		if share == path {
			u.Shares = append(u.Shares[:i], u.Shares[i+1:]...)
			break
		}
	}
}

// Delete the user from the database
func (u *User) Delete(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketUsers)
		return bucket.Delete([]byte(u.Username))
	})
}
