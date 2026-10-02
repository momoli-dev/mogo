// Package mogojwt provides opinionated and ready to use JWT utilities.
package mogojwt

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Issuer creates new JWT using ed25519 keys.
type Issuer struct {
	key ed25519.PrivateKey
}

// NewIssuer creates a new Issuer with the provided key.
func NewIssuer(key ed25519.PrivateKey) *Issuer {
	return &Issuer{key: key}
}

// Issue creates a new JWT token signed with the provided key.
func (i *Issuer) Issue(c jwt.Claims) (string, error) {
	exp, err := c.GetExpirationTime()
	if err != nil || exp == nil || exp.IsZero() {
		return "", errors.New("expiration time can't be empty")
	}

	iss, err := c.GetIssuer()
	if err != nil || iss == "" {
		return "", errors.New("issuer can't be empty")
	}

	nbf, err := c.GetNotBefore()
	if err != nil || nbf == nil || nbf.IsZero() {
		return "", errors.New("not before can't be empty")
	}

	sub, err := c.GetSubject()
	if err != nil || sub == "" {
		return "", errors.New("subject can't be empty")
	}

	return newToken(c, i.key)
}

func newToken(c jwt.Claims, key ed25519.PrivateKey) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, c)
	st, err := t.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return st, nil
}
