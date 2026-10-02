# Agent guidelines for working on `mogo`

`mogo` is a collection of opinionated packages that are used for building Go
systems.

- Use root `Makefile` to run project commands. All recipes are `PHONY`.

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
- `example`: Example `main.go` files showcasing each package.
