NAME := palindrome

build:
	@go build -o $(NAME) ./cmd/server/main.go

clean:
	@rm -f $(NAME)

lint:
	@go get -u github.com/alecthomas/gometalinter
	@gometalinter --install
	@gometalinter --enable-all ./cmd/... ./internal/... ./pkg/...

test-unit:
	@go test -count=1 -race -cover ./...

.PHONY: build clean lint test-unit
