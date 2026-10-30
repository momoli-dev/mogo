# Agent Guidelines for `mogo` Project

Follow these guidelines while working on `mogo` project. `mogo` is a collection
of opinionated packages that are used for building Go systems.

- Never create commits, stage changes, or push to remote repositories.

## Overview

- `mogoapp`: Common application patterns and helpers.
- `mogochi`: Wrapper around chi mux.
- `mogogeo`: Utilities for working with geolocation.
- `mogohash`: Utilities for working with hashes.
- `mogohttp`: Production ready HTTP server.
- `mogohuma`: Opinionated wrapper around `Huma` library with hard conventions
  and easier endpoint registration.
- `mogojwt`: Opinionated JWT wrappers.
- `mogoredis`: Helpers for working with Redis databases.
- `mogoredis/mogoredistest`: Use real Redis database in tests.
- `mogosql`: Helpers for working with SQL databases.
- `mogosql/mogosqltest`: Use real SQL database in tests.
- `mogotest`: Common testing utilities.
- `example`: Runnable applications and tests demonstrating several packages.

## Project Commands

Run project wide commands through the root Makefile:

- `make lint`: Run `golangci-lint` using `.golangci.yaml`.
- `make test`: Start PostgreSQL and Redis with Docker Compose, then run all
  tests with the race detector and coverage.
- `make test-debug`: Run verbose tests with one parallel test at a time.
  Requires PostgreSQL and Redis to be running already.
- `make tidy`: Update module dependencies when dependency changes require it.

## Rules

- Trust parameters supplied by application code; do not add defensive nil,
  empty, or range checks. Validate unknown input from outside the application at
  the boundary where it enters.
- Wrap errors with `%w` to identify the internal operation that failed. Name the
  failing call rather than the enclosing function the caller already knows it
  invoked. Direct forwarding helpers can return an error unchanged when they
  have no useful context to add. For example:

```go
func Foo() error {
	if err := someInternalCall(); err != nil {
		return fmt.Errorf("failed to someInternalCall: %w", err)
	}

	return nil
}
```

## Style Guide

- Keep wrapping messages concise and use `failed to <call>: %w`, for example
  `failed to pool.Begin: %w`. This convention applies to wrapped errors,
  including panics; validation messages and public HTTP errors describe the
  problem itself.
- Write concise Go doc comments for exported symbols and packages, as required
  by the lint configuration. Include caller relevant requirements, ownership,
  defaults, or error behavior when useful.
- Use implementation comments to explain non-obvious intent, constraints, or
  tradeoffs. Avoid narrating visible code in library implementations.
- In examples, keep explanatory comments that guide readers through setup and
  usage, even when they describe visible steps. Examples serve as tutorials.
- Prefer shorter or longer names in relation with the scope of the variable.
  e.g: uid instead of userID, sid instead of sessionID etc. if the scope of
  variable is small enough.
- Group related error handling, function calls and returns together. Separate
  code with empty lines to distinguish groups that do different things. This can
  affect all code, including tests too. e.g:

```go
func Foo(p Product) error {
	keys := []string{
		model.FlowRedisKey(p.Num),
	}
	args := []any{
		p.Code,
		p.FlowID.String(),
	}
	res, err := script.Run(ctx, s.rdb, keys, args...).Text()
	if err != nil {
		return fmt.Errorf("failed to script.Run: %w", err)
	}

	if err := useRes(res); err != nil {
		return fmt.Errorf("failed to useRes: %w", err)
	}

	return nil
}
```
