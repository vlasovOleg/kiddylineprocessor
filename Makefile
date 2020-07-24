.PHONY: tests
tests:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := tests