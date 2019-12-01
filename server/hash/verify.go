package hash

import (
	"crypto/subtle"
	"golang.org/x/crypto/argon2"
)

// Compare a password and generated Argon2id hash
func Verify(password, encodedHash string) (match bool, err error) {
	// Retrieve parameters, salt, and hash
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive key from provided password with same parameters
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Compare hashes
	// Using subtle.ConstantTimeCompare to mitigate timing attacks
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}
