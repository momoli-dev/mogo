package mogojwt

import (
	"crypto/ed25519"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// Validator validates JWT with ed25519 keys.
type Validator struct {
	key ed25519.PublicKey
}

// NewValidator creates a new Validator with the provided key.
func NewValidator(key ed25519.PublicKey) *Validator {
	return &Validator{key: key}
}

// Validate validates the provided JWT token and returns the claims if the token is valid.
//
//nolint:ireturn // Claims can be a custom struct embedding Claims
func (v *Validator) Validate(t string, c jwt.Claims) (jwt.Claims, error) {
	parsed, err := parseToken(t, c, v.key)
	if err != nil {
		return nil, err
	}

	return parsed.Claims, nil
}

// ValidateExpired is like Validate but permits expiration. All other validation failures are returned.
//
//nolint:ireturn // Claims can be a custom struct embedding Claims
func (v *Validator) ValidateExpired(t string, c jwt.Claims) (jwt.Claims, error) {
	parsed, err := parseToken(t, c, v.key)
	if err != nil && (!errors.Is(err, jwt.ErrTokenExpired) || !onlyExpirationErrors(err)) {
		return nil, err
	}

	return parsed.Claims, nil
}

// JWT validation joins failures under ErrTokenInvalidClaims. Inspect every branch
// so expiration cannot hide another standard or custom claims validation failure.
//
//nolint:errorlint // Traverse each immediate child and compare leaves exactly; Is/As could hide sibling failures.
func onlyExpirationErrors(err error) bool {
	switch wrapped := err.(type) {
	case interface{ Unwrap() []error }:
		for _, child := range wrapped.Unwrap() {
			if !onlyExpirationErrors(child) {
				return false
			}
		}
		return true
	case interface{ Unwrap() error }:
		return onlyExpirationErrors(wrapped.Unwrap())
	default:
		return err == jwt.ErrTokenExpired || err == jwt.ErrTokenInvalidClaims
	}
}

func parseToken(t string, c jwt.Claims, key ed25519.PublicKey) (*jwt.Token, error) {
	return jwt.ParseWithClaims(t, c, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodEdDSA.Alg() {
			return nil, errors.New("algorithm not supported")
		}

		return key, nil
	})
}
