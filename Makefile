NAME := palindrome
VERSION ?= latest

build:
	@go build -o $(NAME) ./cmd/server/main.go

build-docker:
	@docker build --build-arg name=$(NAME) --tag $(NAME):$(VERSION) .

lint:
	@gometalinter --enable-all ./cmd/... ./internal/... ./pkg/...

test-unit:
	@go test -count=1 -race -cover ./...

.PHONY: build build-docker lint test-unit
