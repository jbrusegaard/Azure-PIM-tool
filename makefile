.PHONY: build
build:
	@echo "Building the project..."
	@go build -v ./...
	@echo "Build complete."

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests completed."

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@go clean -cache -testcache -modcache
	@echo "Cleanup complete."