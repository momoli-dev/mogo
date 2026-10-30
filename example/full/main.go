// Package main is an example of using all the mogo toolkit for a Go app.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/momoli-dev/mogo/mogoapp"
	"github.com/momoli-dev/mogo/mogochi"
	"github.com/momoli-dev/mogo/mogohttp"
	"github.com/momoli-dev/mogo/mogohuma"
	"github.com/momoli-dev/mogo/mogoredis"
	"github.com/momoli-dev/mogo/mogosql"
	"github.com/sethvargo/go-envconfig/v2"
)

// Config is the app specific configuration sourced from environment variables.
type Config struct {
	ListenAddr      string `env:"LISTEN_ADDR"`
	PostgresConnStr string `env:"POSTGRES_CONN_STR"`
	RedisConnStr    string `env:"REDIS_CONN_STR"`
}

func main() {
	if err := run(); err != nil {
		slog.Error("Application failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	slog.Info("Starting application")

	ctx, cancel := mogoapp.NewAppCtx()
	defer cancel()

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return fmt.Errorf("process config: %w", err)
	}

	// Create and configure your mux.
	chiCfg := mogochi.DefaultConfig()
	chiMux := mogochi.NewChi(chiCfg)

	// Create and configure OpenAPI registry.
	humaRegistry := mogohuma.NewRegistry(
		&mogohuma.NewRegistryParams{
			Adapter:  humachi.NewAdapter(chiMux),
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

	// Create and configure HTTP server.
	httpCfg := mogohttp.DefaultConfig(c.ListenAddr, chiMux)

	pool, err := mogosql.NewPostgresPool(ctx, mogosql.PostgresPoolParams{
		ConnStr:    c.PostgresConnStr,
		HasPostgis: false,
	})
	if err != nil {
		return fmt.Errorf("failed to mogosql.NewPostgresPool: %w", err)
	}
	defer pool.Close()

	rdb, err := mogoredis.NewClient(ctx, c.RedisConnStr)
	if err != nil {
		return fmt.Errorf("create Redis client: %w", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			slog.Error("Failed to close Redis", "err", err)
		}
	}()

	mount(humaRegistry)

	if err := mogohttp.Serve(ctx, httpCfg); err != nil {
		return fmt.Errorf("serve HTTP: %w", err)
	}

	slog.Info("Application stopped")
	return nil
}

// Mount your endpoints here.
func mount(r *mogohuma.Registry) {
	mogohuma.HandleRPC(r, dummy, "auth-phone-start-flow")
	mogohuma.HandleRPC(r, dummy, "auth-phone-finish-flow")
}

func dummy(context.Context, *struct{}) (*struct{}, error) {
	return &struct{}{}, nil
}
