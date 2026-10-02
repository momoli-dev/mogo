// Package main is an example of mogohuma.
package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/momoli-dev/mogo/mogohuma"
)

func main() {
	// Create your mux.
	mux := http.NewServeMux()

	// Create a registry.
	r := mogohuma.NewRegistry(
		&mogohuma.NewRegistryParams{
			Adapter:  humago.NewAdapter(mux, ""),
			Title:    "Petstore",
			Version:  "1.0.0",
			DocsPath: "/swagger",
			SecureMiddlewareMaker: func(_ huma.API) mogohuma.Middleware {
				return func(ctx huma.Context, next func(huma.Context)) {
					slog.Warn("Secure middleware is not implemented. All requests are allowed.")
					next(ctx)
				}
			},
		},
	)

	// Register your endpoints, fluently and with options.
	mogohuma.Get(r, "/pets", dummyHandler, "list-pets", mogohuma.WithInsecure())
	mogohuma.Get(r, "/pets/{id}", dummyHandler, "get-pet", mogohuma.WithInsecure())
	mogohuma.Post(r, "/pets", dummyHandler, "create-pet", mogohuma.WithMiddlewares(dummyMiddleware))
	mogohuma.Put(r, "/pets/{id}", dummyHandler, "replace-pet")
	mogohuma.Patch(r, "/pets/{id}", dummyHandler, "rename-pet")
	mogohuma.Delete(r, "/pets/{id}", dummyHandler, "delete-pet")
	mogohuma.Get(r, "/inventory", dummyHandler, "inventory")

	// You may also register some opinionated RPC endpoints.
	mogohuma.HandleRPC(r, dummyHandler, "adopt-pet")

	// Export the automatically generated spec without starting a server.
	if err := r.CreateJSONSpecFile("petstore.openapi.json"); err != nil {
		panic(err)
	}
	if err := r.CreateYAMLSpecFile("petstore.openapi.yaml"); err != nil {
		panic(err)
	}
	slog.Info("OpenAPI exported", "title", r.HumaAPI.OpenAPI().Info.Title)
}

func dummyHandler(context.Context, *struct{}) (*struct{}, error) {
	return &struct{}{}, nil
}

func dummyMiddleware(ctx huma.Context, next func(huma.Context)) {
	ctx.SetHeader("X-Petstore", "example")
	next(ctx)
}
