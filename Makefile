.PHONY: lint
lint: ## Performs linting
	golangci-lint run -c ./build/.golangci-lint.yml ./...

.PHONY: unit_test
unit_test: ## Performs unit testing with minimal output
	go test --tags="unit_test" ./...

.PHONY: coverage
coverage_file="/tmp/go-cover.$$.tmp"
coverage: ## Creates a test coverage report
	go test --tags="unit" -coverprofile=$coverage_file ./... && go tool cover -html=$coverage_file && unlink $coverage_file

.PHONY: build
build: ## Compiles the server
	go build -o target/trackpal ./cmd/trackpal/...

help: ## Displays this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "};\
	 {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
