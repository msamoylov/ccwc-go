.PHONY: all
all: lint test

.PHONY: lint
lint: golangci-lint tidy-lint

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: tidy-lint
tidy-lint:
	go mod tidy && \
	git diff --exit-code -- go.mod go.sum

.PHONY: test
test:
	go test -race ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./... && \
	go tool cover -html=cover.out -o cover.html

.PHONY: build
build:
	 go build ./cmd/ccwc