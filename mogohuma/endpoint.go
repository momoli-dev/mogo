// Package mogohuma is an opinionated wrapper around Huma for creating HTTP APIs.
package mogohuma

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// Endpoint defines a route and its middleware.
type Endpoint struct {
	// Name is the unique OpenAPI operation ID. This will be used by code generators.
	Name string

	// Secure enables [Registry.SecureMiddleware] and OpenAPI bearer auth metadata.
	Secure bool

	Method      string
	Path        string
	Middlewares huma.Middlewares
}

// Handler handles a typed request and response.
type Handler[TReq, TRes any] func(context.Context, *TReq) (*TRes, error)

// Middleware intercepts a Huma request.
type Middleware func(ctx huma.Context, next func(huma.Context))

// MiddlewareMaker creates middleware for a Huma API.
type MiddlewareMaker func(humaAPI huma.API) Middleware
