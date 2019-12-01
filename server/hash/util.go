package hash

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

var (
	ErrInvalidHash         = errors.New("encoded hash is not in correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

// Argon2id configuration parameters
type Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// Generate some number of cryptographically random bytes
func randomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// Encode to standard hash format
func encodeHash(hash, salt []byte, p *Params) string {
	// Convert hash and salt to base64
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	// Return formatted string
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)
}

// Decode from standard hash form
func decodeHash(encodedHash string) (p *Params, salt, hash []byte, err error) {
	// Ensure correct number of segments
	segments := strings.Split(encodedHash, "$")
	if len(segments) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	// Retrieve and validate argon2id version
	var version int
	if _, err := fmt.Sscanf(segments[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	} else if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Load config parameters
	p = &Params{}
	if _, err := fmt.Sscanf(segments[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism); err != nil {
		return nil, nil, nil, err
	}

	// Decode salt and set salt length
	salt, err = base64.RawStdEncoding.DecodeString(segments[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	// Decode hash and set key length
	hash, err = base64.RawStdEncoding.DecodeString(segments[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return
}
