FROM golang:alpine AS builder
ARG name
WORKDIR /go/src/github.com/nicholaslam/palindrome-service/
COPY . /go/src/github.com/nicholaslam/palindrome-service/
RUN go build -o $name ./cmd/server/main.go

FROM alpine:latest
ARG name
ENV NAME=$name
COPY --from=builder /go/src/github.com/nicholaslam/palindrome-service/$name /usr/local/bin/
# Run as non-root user.
RUN chown -R nobody:nogroup /usr/local/bin/$name && chmod +x /usr/local/bin/$name
USER nobody
# Cannot use ARG in ENTRYPOINT.
ENTRYPOINT $NAME
EXPOSE 8080
