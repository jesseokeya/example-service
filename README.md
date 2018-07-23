# palindrome-service

Palindrome service allows create, read, read all, and delete operations on messages. Messages are evaluated to determine if they are palindromes. Only English is supported.

## Required Software

- Use [dep](https://github.com/golang/dep) for package management.
- Use [gometalinter](https://github.com/alecthomas/gometalinter) for linting.

## Building and Running

Use `make build` to build the service. Execute the `palindrome` binary to start the service.

## Testing

Use `make test-unit` to run the unit tests.

## Linting

Use `make lint` to lint the source code.
