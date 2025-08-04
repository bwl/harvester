# Simple Makefile for Harvest of Stars

.PHONY: run build fmt imports vet lint staticcheck test clean watch-new watch-modified watch-go

run:
	go run ./cmd/game

build:
	go build ./...

fmt:
	go fmt ./...

imports:
	@command -v goimports >/dev/null 2>&1 && goimports -w . || echo "goimports not installed; skipping"

vet:
	go vet ./...

staticcheck:
	@command -v staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "staticcheck not installed; skipping"

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed; skipping"

test:
	go test ./...

clean:
	go clean -testcache

# File watching targets
watch-new:
	./scripts/watch-new-files.sh

watch-modified:
	./scripts/watch-modified-files.sh

watch-go:
	./scripts/watch-go-files.sh
