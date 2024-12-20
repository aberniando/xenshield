include .env.example
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker compose up --build postgres app
.PHONY: compose-up

compose-up-integration-test: ### Run docker-compose with integration test
	docker compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down: ### Down docker-compose
	docker compose down --remove-orphans
.PHONY: compose-down

lint: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

test:
	go run github.com/onsi/ginkgo/v2/ginkgo -r -p --succinct --race --trace --keep-going --cover --coverprofile=$(COVER_PROFILE_FILE) --timeout=20m --json-report=report.json --coverpkg=github.com/aberniando/xenshield/internal/usecases/...,github.com/aberniando/xenshield/internal/handler/... --skip-package=integration

NAME ?= unnamed_migration

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations $(NAME)
.PHONY: migrate-create
