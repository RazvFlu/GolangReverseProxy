# Go HTTP reverse proxy

This is a simple HTTP reverse proxy written in Go.

Implemented features:
1. The reverse proxy should be able to forward any kind of request but inspect only GET
requests.
2. The proxy should inspect the response body for any sensitive information and mask it
before forwarding the response.
3. The proxy should log all incoming requests, including headers and body, and response
headers and body.
4. The proxy should be able to block requests based on a set of predefined rules. For
example, you can block requests that contain specific headers or parameters.
5. You should provide a configuration file that allows the startup to define the rules for
blocking requests

## First time setup

Install golang: https://golang.org/doc/install
Install dependencies. In the root folder of the project run:

```bash
go get ./...
```

## How to run

```bash
go run ./cmd/reverseproxyserver/main.go
```

The server loads the configuration from ./internal/config/test_data.json
By default it will start a server which listens on localhost:8080 and will proxy requests based on the configured paths.

In this example the code will start some mirror backend services on localhost:10001 and localhost:10002 for testing purposes.

The configured paths in this example are: test1 and test2
test1 will forward the requests towards http://localhost:10001
test2 will forward the requests towards http://localhost:10002

test1 is configured to block requests which contain the parameter "password"
test2 is configured to block requests which contain the header "Auth"

## Logging

```bash
tail -f ./logs/application.log
```

## Testing

Unit tests:

```bash
go test ./...
```

Manual testing:

Test GET request with sensitive data, without any blocking.
The response should contain the masked email address.

```bash
curl -X GET -H "Content-Type: text/plain" -d 'This is an example email: john.doe@example.com' http://localhost:8080/test1 -v
```

Test GET request with sensitive data, with query param blocking:
The returned status code will be 403 Forbidden.

```bash
curl -X GET -H "Content-Type: text/plain" -d 'This is an example email: john.doe@example.com' http://localhost:8080/test1?password=test -v
```

Test GET request with sensitive data, with header blocking.
The returned status code will be 403 Forbidden.

```bash
curl -X GET -H "Content-Type: text/plain" -H "Auth: YourTokenHere" -d 'This is an example email: john.doe@example.com' http://localhost:8080/test2 -v
```

Test POST request with sensitive data which should not be inspected.
The response should contain the original email address and status code 200 OK.

```bash
curl -X POST -H "Content-Type: text/plain" -d 'This is an example email: john.doe@example.com' http://localhost:8080/test1 -v
```
# GolangReverseProxy
