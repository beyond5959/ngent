.PHONY: test run fmt build-web build

test:
	go test ./...

build-web:
	cd internal/webui/web && npm ci && npm run build

build: build-web
	go build -o bin/ngent ./cmd/ngent

run: build-web
	go run ./cmd/ngent

run-local: build-web
	go run ./cmd/ngent --allow-public=false

fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')
