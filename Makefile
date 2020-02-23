.PHONY: lint
lint: ## Performs linting
	golangci-lint run -c ./build/.golangci-lint.yml ./...

.PHONY: unit_test
unit_test: ## Performs unit testing with minimal output
	go test --tags="unit_test" ./...

.PHONY: integration_test
integration_test: ## Perform integration testing (requires env: POSTGRES_DSN)
	go test --tags="integration_test" ./...

.PHONY: coverage
coverage_file=./output/coverage.cov
coverage: ## Creates a test coverage report (unit tests)
	mkdir -p output
	go test --tags="unit_test" -coverprofile=$coverage_file ./... && go tool cover -func=$coverage_file

.PHONY: coverage_full
coverage_file=./output/coverage.cov
coverage_full: ## Creates a test coverage report (requires env: POSTGRES_DSN)
	mkdir -p output
	go test --tags="unit_test integration_test" -coverprofile=${coverage_file} ./... && go tool cover -func=${coverage_file}

.PHONY: build
build: ## Compiles the server
	go build -o target/trackpal ./cmd/trackpal/...

.PHONY: deploy
deploy: ## Deploys the built artifacts
	scp target/trackpal aws.trackpal:trackpal

help: ## Displays this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "};\
	 {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
