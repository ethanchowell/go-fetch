
VERSION ?= $(shell git describe --always --dirty --tags 2>/dev/null || echo "undefined")
# Allow to override image registry.
DOCKER_REGISTRY ?= quay.io
DOCKER_NAMESPACE ?= ethanchowell
PROJECT_NAME ?= go-fetch

IMAGE ?= $(DOCKER_REGISTRY)/$(DOCKER_NAMESPACE)/$(PROJECT_NAME):$(VERSION)

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Housekeeping

.PHONY: fmt
fmt: ## Run go fmt ./...
	@go fmt ./...

.PHONY: vet
vet: fmt ## Run go vet ./...
	@go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	@golangci-lint run

##@ Building

.PHONY: build
build: fmt vet lint ## Build the executable
	@CGO_ENABLED=0 GOARCH=amd64 go build -buildvcs=false -ldflags="-w -s -X version.VERSION=${VERSION}" -o ./build/go-fetch ./cmd/go-fetch

.PHONY: docker
docker: build ## Build the docker image
	@docker build -t $(IMAGE) -f Dockerfile .