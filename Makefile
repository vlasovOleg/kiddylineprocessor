.PHONY: tests
tests:
	go test -v -race -timeout 30s ./...

.PHONY: lint
lint:
	@docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v

.DEFAULT_GOAL := tests