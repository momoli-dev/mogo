package mogohuma

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

const (
	bearerScheme = "bearer"
)

// Registry holds all the endpoints and configuration for a Huma API.
type Registry struct {
	HumaAPI huma.API

	// SecureMiddleware runs before endpoint middleware on secure routes.
	SecureMiddleware Middleware
}

// NewRegistryParams configures [NewRegistry].
type NewRegistryParams struct {
	// Adapter connects Huma to the application's router.
	Adapter huma.Adapter

	Title                 string
	Version               string
	SecureMiddlewareMaker MiddlewareMaker

	// DocsPath is where the docs UI is served from. If empty, UI is not served.
	DocsPath string
}

// NewRegistry creates a Huma API with OpenAPI, documentation, and bearer auth metadata.
// Adapter and SecureMiddlewareMaker must be non-nil.
func NewRegistry(p *NewRegistryParams) *Registry {
	humaConfig := huma.DefaultConfig(p.Title, p.Version)
	humaConfig.DocsPath = p.DocsPath
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		bearerScheme: {
			Type:   "http",
			Scheme: bearerScheme,
		},
	}
	humaAPI := huma.NewAPI(humaConfig, p.Adapter)

	secureMiddleware := p.SecureMiddlewareMaker(humaAPI)

	return &Registry{
		HumaAPI:          humaAPI,
		SecureMiddleware: secureMiddleware,
	}
}

// Register adds an endpoint after applying opts.
func Register[TReq, TRes any](r *Registry, pe *Endpoint, h Handler[TReq, TRes], opts ...Option) {
	e := *pe
	for _, opt := range opts {
		opt(&e)
	}

	middlewares := huma.Middlewares{}
	securitySchema := []map[string][]string{}

	if e.Secure {
		middlewares = append(middlewares, r.SecureMiddleware)
		securitySchema = []map[string][]string{
			{bearerScheme: {}},
		}
	}

	middlewares = append(middlewares, e.Middlewares...)

	//nolint:exhaustruct // [huma.Operation]
	huma.Register(r.HumaAPI, huma.Operation{
		OperationID:   e.Name,
		Method:        e.Method,
		Path:          e.Path,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
		Security:      securitySchema,
	}, h)
}
