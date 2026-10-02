.PHONY: *

all: tidy lint test

tidy:
	go mod tidy

lint:
	golangci-lint run

test:
	docker compose -f compose.test.yaml up -d --wait
	go test -coverprofile=coverage.out -race ./...

test-debug:
	go test -v -parallel=1 ./...

clean:
	go clean -testcache
	golangci-lint cache clean
	docker compose -f compose.test.yaml down --volumes

doc:
	pkgsite -open
