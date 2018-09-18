FROM golang:alpine AS builder
WORKDIR /go/src/github.com/nicholaslam/example-service/
COPY . /go/src/github.com/nicholaslam/example-service/
RUN go build -o palindrome ./cmd/server/main.go

FROM alpine:latest
COPY --from=builder /go/src/github.com/nicholaslam/example-service/palindrome /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/palindrome && chmod +x /usr/local/bin/palindrome
USER nobody
ENTRYPOINT palindrome
EXPOSE 8080
