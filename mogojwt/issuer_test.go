package mogojwt_test

import (
	"crypto/ed25519"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/momoli-dev/mogo/mogojwt"
)

func testClaims(t *testing.T) jwt.RegisteredClaims {
	t.Helper()

	return jwt.RegisteredClaims{
		Issuer:    "test-iss",
		Subject:   "test-sub",
		Audience:  []string{},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	}
}

func TestIssuer_Issue(t *testing.T) {
	t.Parallel()
	_, key, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	issuer := mogojwt.NewIssuer(key)
	require.NotNil(t, issuer)

	claims := testClaims(t)

	token, err := issuer.Issue(claims)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	has3Parts := strings.Count(token, ".") == 2
	require.True(t, has3Parts)
}

func TestIssuer_IssueMalformedClaims(t *testing.T) {
	t.Parallel()
	_, key, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	issuer := mogojwt.NewIssuer(key)
	require.NotNil(t, issuer)

	okClaims := testClaims(t)

	c := okClaims
	c.ExpiresAt = jwt.NewNumericDate(time.Time{})
	_, err = issuer.Issue(c)
	require.Error(t, err)

	c = okClaims
	c.Issuer = ""
	_, err = issuer.Issue(c)
	require.Error(t, err)

	c = okClaims
	c.NotBefore = jwt.NewNumericDate(time.Time{})
	_, err = issuer.Issue(c)
	require.Error(t, err)

	c = okClaims
	c.Subject = ""
	_, err = issuer.Issue(c)
	require.Error(t, err)
}
