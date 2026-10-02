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

// NewHashedString creates a new hash of passed in string.
func NewHashedString(s string) HashedString {
	hash, err := argon2id.CreateHash(s, argon2id.DefaultParams)
	if err != nil {
		panic(fmt.Errorf("mogohash failed to create hash: %w", err))
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
