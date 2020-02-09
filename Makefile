.PHONY: lint
lint: ## Performs linting
	golangci-lint run -c ./build/.golangci-lint.yml ./...

.PHONY: docs
docs: ## Compiles the docs
	mkdir -p ./target/docs
	docsify init ./target/docs

.PHONY: unit_test
unit_test: ## Performs unit testing with minimal output
	go test --tags="unit_test" ./...

help: ## Displays this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "};\
	 {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
