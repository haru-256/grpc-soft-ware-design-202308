.DEFAULT_GOAL := help

.PHONY: init
init: ## Initial setup
	mise install
	buf config init

.PHONY: generate
generate:  ## Generate gRPC code
	buf generate

.PHONY: lint
lint:  ## Lint proto and go files
	buf lint
	golangci-lint run --config=./.golangci.yml ./...

.PHONY: fmt
fmt:  ## Format proto and go files
	buf format -w .
	go fmt ./...

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
