// Package mogohash provides common hashing utilities.
package mogohash

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

// HashedString represents a hashed and secured string.
type HashedString struct {
	hash string
}

// NewHashedString hashes s using the default Argon2id parameters and panics on error.
func NewHashedString(s string) HashedString {
	hash, err := argon2id.CreateHash(s, argon2id.DefaultParams)
	if err != nil {
		panic(fmt.Errorf("failed to argon2id.CreateHash: %w", err))
	}

	return HashedString{
		hash: hash,
	}
}

// Check checks if a string matches the hash.
func (p *HashedString) Check(s string) bool {
	match, err := argon2id.ComparePasswordAndHash(s, p.hash)
	if err != nil {
		return false
	}

	return match
}
