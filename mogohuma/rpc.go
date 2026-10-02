package mogohuma

import (
	"context"
	"errors"
)

// RPCRequest wraps an RPC input as a Huma request body.
type RPCRequest[T any] struct {
	Body *T
}

// RPCResponse wraps an RPC result as a Huma response body.
type RPCResponse[T any] struct {
	Body *T
}

// RPCFailure is a public error message returned as HTTP 400, even when wrapped.
//
//nolint:errname // Failure is the RPC domain term and part of the public API.
type RPCFailure string

// Error returns the public message.
func (f RPCFailure) Error() string {
	return string(f)
}

func rpcErrToHTTPError(err error) error {
	if err == nil {
		return nil
	}

	// Expose only the failure message, not wrapping errors.
	failure, ok := errors.AsType[RPCFailure](err)
	if ok {
		return Err400(failure)
	}

	return Err500(err)
}

// MakeRPCHandler adapts a handler to Huma request and response bodies.
// An [RPCFailure] becomes HTTP 400; other errors are logged and hidden behind HTTP 500.
func MakeRPCHandler[TReq, TRes any](rpc Handler[TReq, TRes]) Handler[RPCRequest[TReq], RPCResponse[TRes]] {
	return func(ctx context.Context, req *RPCRequest[TReq]) (*RPCResponse[TRes], error) {
		res, err := rpc(ctx, req.Body)
		if err != nil {
			return nil, rpcErrToHTTPError(err)
		}

		if res == nil {
			return nil, nil
		}

		return &RPCResponse[TRes]{
			Body: res,
		}, nil
	}
}
