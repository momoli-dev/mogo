package mogohttp_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/momoli-dev/mogo/mogohttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeReturnsListenError(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	listener, err := new(net.ListenConfig).Listen(ctx, "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, listener.Close())
	})

	err = mogohttp.Serve(
		ctx,
		mogohttp.DefaultConfig(listener.Addr().String(), http.NotFoundHandler()),
	)
	require.Error(t, err)
	assert.ErrorContains(t, err, "failed to s.ListenAndServe")
}

func TestServeShutsDownWhenContextIsCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	result := make(chan error, 1)
	go func() {
		result <- mogohttp.Serve(
			ctx,
			mogohttp.DefaultConfig("127.0.0.1:0", http.NotFoundHandler()),
		)
	}()

	select {
	case err := <-result:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("ServeHTTP did not return after cancellation")
	}
}
