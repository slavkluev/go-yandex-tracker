.PHONY: fmt lint test coverage check

## fmt: format code with goimports
fmt:
	go tool goimports -w .

## lint: run golangci-lint (matches CI lint job)
lint:
	golangci-lint run

## test: run tests with race detector (matches CI test job)
test:
	go test -race ./...

## coverage: generate coverage profile
coverage:
	go test -coverprofile=coverage.out ./...

## check: run lint + test (pre-push validation)
check: lint test
