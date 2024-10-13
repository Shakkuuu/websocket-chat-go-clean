TEST_FLAGS := -v -cover -timeout 30s

.PHONY: test
test:
	go test $(TEST_FLAGS) ./...

.PHONY: lint
lint:
	golangci-lint run ./...
