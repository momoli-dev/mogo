package mogojwt_test

import (
	"crypto/ed25519"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/golang-jwt/jwt/v5"
	"github.com/momoli-dev/mogo/mogojwt"
	"github.com/momoli-dev/mogo/mogotest"
)

func TestValidator_Validate(t *testing.T) {
	t.Parallel()
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	issuer := mogojwt.NewIssuer(privKey)
	require.NotNil(t, issuer)

	claims := jwt.RegisteredClaims{
		Subject:   "test-sub",
		Issuer:    "test-iss",
		Audience:  []string{},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	}
	token, err := issuer.Issue(claims)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	validator := mogojwt.NewValidator(pubKey)
	require.NotNil(t, validator)

	res, err := validator.Validate(token, &jwt.RegisteredClaims{})
	require.NoError(t, err)
	require.NotNil(t, res)

	resClaims, ok := res.(*jwt.RegisteredClaims)
	require.True(t, ok)

	// NOTE: Other fields may have been altered by the validator, so we only check the ones we care about.
	require.Equal(t, claims.Subject, resClaims.Subject)
	require.Equal(t, claims.Issuer, resClaims.Issuer)
	require.NotNil(t, resClaims.ExpiresAt)
	mogotest.RequireEqualTime(t, claims.ExpiresAt.Time, resClaims.ExpiresAt.Time)
	require.Equal(t, claims.ID, resClaims.ID)
}

func TestValidator_ValidateNoAlg(t *testing.T) {
	t.Parallel()

	//nolint:lll // Tokens are long.
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTc2NTExNjc3NCwiZXhwIjoxNzY1MTIwMzc0fQ"

	key, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	validator := mogojwt.NewValidator(key)
	require.NotNil(t, validator)

	_, err = validator.Validate(token, &jwt.RegisteredClaims{})
	require.Error(t, err)
}

func TestValidator_ValidateExpiredNoAlg(t *testing.T) {
	t.Parallel()

	//nolint:lll // Tokens are long.
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTc2NTExNjc3NCwiZXhwIjoxNzY1MTIwMzc0fQ"

	key, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	validator := mogojwt.NewValidator(key)
	require.NotNil(t, validator)

	_, err = validator.ValidateExpired(token, &jwt.RegisteredClaims{})
	require.Error(t, err)
}

func TestValidator_ValidateExpiredWrongAlg(t *testing.T) {
	t.Parallel()

	//nolint:lll // Tokens are long.
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTc2NTExNjc3NCwiZXhwIjoxNzY1MTIwMzc0fQ.rElgPn3ZGG7QvHmv7Zr69NOuRIOrXJxGGbxEIxUUv5k"

	key, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	validator := mogojwt.NewValidator(key)
	require.NotNil(t, validator)

	//nolint:exhaustruct // struct passed for type selection
	_, err = validator.ValidateExpired(token, &jwt.RegisteredClaims{})
	require.Error(t, err)
}

func TestValidator_ValidateExpiredOK(t *testing.T) {
	t.Parallel()
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	issuer := mogojwt.NewIssuer(privKey)
	require.NotNil(t, issuer)

	claims := testClaims(t)
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
	token, err := issuer.Issue(claims)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	validator := mogojwt.NewValidator(pubKey)
	require.NotNil(t, validator)

	//nolint:exhaustruct // struct passed for type selection
	res, err := validator.ValidateExpired(token, &jwt.RegisteredClaims{})
	require.NoError(t, err)
	require.NotNil(t, res)

	resClaims, ok := res.(*jwt.RegisteredClaims)
	require.True(t, ok)

	// NOTE: Other fields may have been altered by the validator, so we only check the ones we care about.
	require.Equal(t, claims.Subject, resClaims.Subject)
	require.Equal(t, claims.Issuer, resClaims.Issuer)
	require.NotNil(t, resClaims.ExpiresAt)
	mogotest.RequireEqualTime(t, claims.ExpiresAt.Time, resClaims.ExpiresAt.Time)
	require.Equal(t, claims.ID, resClaims.ID)
}

type rejectedClaims struct {
	jwt.RegisteredClaims
}

func (*rejectedClaims) Validate() error {
	return errors.New("custom claims rejected")
}

func TestValidator_ValidateExpiredPreservesOtherFailures(t *testing.T) {
	t.Parallel()
	const expired = "expired"
	pub, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	otherPub, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	for _, name := range []string{expired, "unexpired", "not yet valid", "custom failure", "wrong key"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			claims := testClaims(t)
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
			var decoded jwt.Claims = &jwt.RegisteredClaims{}
			key := pub
			var wantErr error
			switch name {
			case expired:
			case "unexpired":
				claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
			case "not yet valid":
				claims.NotBefore = jwt.NewNumericDate(time.Now().Add(time.Hour))
				wantErr = jwt.ErrTokenNotValidYet
			case "custom failure":
				decoded = &rejectedClaims{RegisteredClaims: jwt.RegisteredClaims{}}
				wantErr = jwt.ErrTokenInvalidClaims
			case "wrong key":
				key = otherPub
				wantErr = jwt.ErrTokenSignatureInvalid
			}
			token, err := mogojwt.NewIssuer(priv).Issue(claims)
			require.NoError(t, err)
			validator := mogojwt.NewValidator(key)
			got, err := validator.ValidateExpired(token, decoded)
			if wantErr != nil {
				require.ErrorIs(t, err, wantErr)
				require.Nil(t, got)

				return
			}
			require.NoError(t, err)
			require.Same(t, decoded, got)
			if name == expired {
				_, err = validator.Validate(token, &jwt.RegisteredClaims{})
				require.ErrorIs(t, err, jwt.ErrTokenExpired)
			}
		})
	}
}
