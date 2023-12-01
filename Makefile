.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: build
build:
	go build -race .

.PHONY: test
test: ## Run unit tests, skipping the ones that require a real database.
	go test -race -v .
