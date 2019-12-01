package hash

import "golang.org/x/crypto/argon2"

// Generate an Argon2id hash of the provided password with the default recommended configuration
func DefaultHash(password string) (string, error) {
	return Hash(password, &Params{
		memory:      32 * 1024,
		iterations:  4,
		parallelism: 4,
		saltLength:  16,
		keyLength:   32,
	})
}

// Generate an Argon2id hash of the provided password with specified configuration
func Hash(password string, p *Params) (encodedHash string, err error) {
	// Generate cryptographically secure salt
	salt, err := randomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	// Generate hash of password
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return encodeHash(hash, salt, p), nil
}
