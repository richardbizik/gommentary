#
#
#
#
#
#
DEFAULT_CONFIG = api

ifeq ("$(shell go env GOOS)","windows")
	GO_PATH := $(subst \,/,$(shell go env GOPATH))
else
	GO_PATH := $(GOPATH)
endif

# process arguments and save to run-args if we are running build-* or run-* target if nothing specified use default
TARGETS_WITH_PARAMS := run-dev run-prod build build-image test-all bench-all test-report build-for-image format-migration
ifneq (, $(filter $(firstword $(MAKECMDGOALS)),$(TARGETS_WITH_PARAMS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
ifeq (,$(RUN_ARGS))
  RUN_ARGS := $(DEFAULT_CONFIG)
endif
endif

##@ Development

.PHONY: generate ## Run all generators
generate: oapi-generate sqlc-generate

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: run-dev
run-dev: ## Start the project with dev profile and configuration
	export PROFILE=DEV; \
	export CONFIG_FILE=$(CURDIR)/conf/$(RUN_ARGS)/conf-dev.yaml; \
	go run cmd/$(RUN_ARGS)/*.go

.PHONY: sqlc
SQLC ?= $(LOCALBIN)/sqlc
sqlc: ## Download open-api generator locally if necessary.
ifeq (,$(wildcard $(SQLC)))
ifeq (,$(shell which sqlc 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(SQLC)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	GOBIN=$(shell pwd)/bin go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest ;\
	chmod +x $(SQLC) ;\
	}
else
SQLC = $(shell which sqlc)
endif
endif

.PHONY: sqlc-generate
sqlc-generate: sqlc ## run sqlc generator
	$(SQLC) generate

.PHONY: oapi
OAPI ?= $(LOCALBIN)/oapi-codegen
oapi: ## Download open-api generator locally if necessary.
ifeq (,$(wildcard $(OAPI)))
ifeq (,$(shell which oapi-codegen 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OAPI)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	GOBIN=$(shell pwd)/bin go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest ;\
	chmod +x $(OAPI) ;\
	}
else
OAPI = $(shell which oapi-codegen)
endif
endif

.PHONY: oapi-generate
oapi-generate: oapi ## Generate rest server files based on open api spec
	$(OAPI) -generate types -package handlers open-api.yaml > ./internal/rest/handlers/zz_generated_types.go
	$(OAPI) -generate chi-server,strict-server -package handlers open-api.yaml > ./internal/rest/handlers/zz_generated_server.go
	$(OAPI) -generate spec -package handlers open-api.yaml > ./internal/rest/handlers/zz_generated_spec.go

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: bench
bench: ## Run benchmarks
	go test ./... -bench=.

##@ Build

.PHONY: docker-build
docker-build: ## Build docker image
	docker build --progress=plain . -t gommentary:latest

##@ Misc

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

