.PHONY: test run fmt build-web build

build: build-web
	go build -o bin/ngent ./cmd/ngent

test:
	go test ./...

build-web:
	cd internal/webui/web && npm ci && npm run build

run: build-web
	go run ./cmd/ngent

run-local: build-web
	go run ./cmd/ngent --allow-public=false

fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')
