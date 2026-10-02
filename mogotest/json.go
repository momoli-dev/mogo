package mogotest

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// JSONUnmarshal decodes JSON into T, failing the test on error.
//
//nolint:ireturn // Returns the requested type.
func JSONUnmarshal[T any](t *testing.T, data []byte) T {
	t.Helper()

	var value T
	err := json.Unmarshal(data, &value)
	require.NoError(t, err)

	return value
}

// JSONUnmarshalResponse decodes response JSON without consuming the body.
//
//nolint:ireturn // Returns the requested type.
func JSONUnmarshalResponse[T any](t *testing.T, res *httptest.ResponseRecorder) T {
	t.Helper()

	return JSONUnmarshal[T](t, res.Body.Bytes())
}
