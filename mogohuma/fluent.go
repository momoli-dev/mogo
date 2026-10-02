package mogohuma

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// Get registers a GET endpoint, secure unless [WithInsecure] is supplied.
func Get[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodGet,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// Post registers a POST endpoint, secure unless [WithInsecure] is supplied.
func Post[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodPost,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// Put registers a PUT endpoint, secure unless [WithInsecure] is supplied.
func Put[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodPut,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// Delete registers a DELETE endpoint, secure unless [WithInsecure] is supplied.
func Delete[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodDelete,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// Patch registers a PATCH endpoint, secure unless [WithInsecure] is supplied.
func Patch[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodPatch,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// HandleRPC registers `POST /rpc/<name>` with HTTP 200 on success.
// Endpoints are secure unless [WithInsecure] is supplied.
func HandleRPC[TReq, TRes any](
	r *Registry,
	handler Handler[TReq, TRes],
	name string,
	opts ...Option,
) {
	rpcHandler := MakeRPCHandler(handler)
	Post(r, "/rpc/"+name, rpcHandler, name, opts...)
}
