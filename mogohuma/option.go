package mogohuma

// Option modifies an endpoint before registration.
type Option func(*Endpoint)

// WithInsecure disables registry authentication and OpenAPI bearer auth metadata.
func WithInsecure() Option {
	return func(e *Endpoint) {
		e.Secure = false
	}
}

// WithMiddlewares appends middleware to the endpoint.
func WithMiddlewares(m ...Middleware) Option {
	return func(e *Endpoint) {
		for _, middleware := range m {
			e.Middlewares = append(e.Middlewares, middleware)
		}
	}
}
