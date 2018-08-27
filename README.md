# palindrome-service

Palindrome service allows create, read, read all, and delete operations on messages. Messages are evaluated to determine if they are palindromes. Only English is supported.

## Required Software

- Use [dep](https://github.com/golang/dep) for package management.
- Use [gometalinter](https://github.com/alecthomas/gometalinter) for linting.

## Building and Running

Use `make build` to build the service. Execute the `palindrome` binary to start the service. The supported command-line flags are `http-addr` and `strict-palindrome`.

```
./palindrome -http-addr=:8081 -strict-palindrome=false
```

Use `make build-docker` to build the docker image. Use `docker run` to start the service in a container. The supported environment variables are `HTTP_ADDR` and `STRICT_PALINDROME`.


```
docker run -e HTTP_ADDR=:8081 -e STRICT_PALINDROME=false -p 8081:8081 palindrome:latest
```

## Testing

Use `make test-unit` to run the unit tests.

## Linting

Use `make lint` to lint the source code.
