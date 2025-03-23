# ==============================================================================
# Define dependencies

SERVICE_VERSION            := 0.1.0
SERVICE_NAME               := user-service
BUILD_VERSION              ?= $(shell git symbolic-ref HEAD 2> /dev/null | cut -b 12-)_$(shell git log --pretty=format:%h -1)
BUILD_DATE                 ?= $(shell date +%FT%T%z)

BASE_IMAGE                 := gravitum
IMAGE_SERVICE              := $(BASE_IMAGE)/service/user:$(SERVICE_VERSION)

MAIN_PACKAGE_PATH          := ./cmd/
BINARY_NAME                := user_service
TEMP_BIN                   := /var/tmp/gravitum/user/bin
PROJECT_PKG                := github.com/Karzoug/gravitum-user-service

LDFLAGS += -s -w
LDFLAGS += -X ${PROJECT_PKG}/pkg/buildinfo.buildVersion=${BUILD_VERSION} -X ${PROJECT_PKG}/pkg/buildinfo.buildDate=${BUILD_DATE} -X ${PROJECT_PKG}/pkg/buildinfo.serviceVersion=$(SERVICE_VERSION)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test fmt lint
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## fmt: format .go files
.PHONY: fmt
fmt:
	go run golang.org/x/tools/cmd/goimports@latest -local=${PROJECT_PKG} -l -w  .
	go run mvdan.cc/gofumpt@latest -l -w  .

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## lint: run linters
.PHONY: lint
lint:
	golangci-lint run ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	go build -ldflags "${LDFLAGS}" -o ${TEMP_BIN}/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## generate: generate all necessary code
.PHONY: generate
generate:
	go generate ./...

## clean: clean all temporary files
.PHONY: clean
clean:
	rm -rf $(TEMP_DIR)

# ==============================================================================
# Install dependencies

## dev-install-deps: install dependencies with fixed versions in a temporary directory
dev-install-deps:
	go install tool

# ==============================================================================
# Building containers

## service: build the service image
.PHONY: service
service:
	docker build \
		-f build/dockerfile.service \
		-t $(IMAGE_SERVICE) \
		--build-arg BUILD_REF=$(BUILD_VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg VERSION=$(SERVICE_VERSION) \
		--build-arg PROJECT_PKG=$(PROJECT_PKG) \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		.

# ------------------------------------------------------------------------------

dev-compose-up:
	IMAGE_SERVICE=$(IMAGE_SERVICE) docker compose -f deploy/docker/dev-compose.yaml  up -d

dev-compose-down:
	IMAGE_SERVICE=$(IMAGE_SERVICE) docker compose -f deploy/docker/dev-compose.yaml down
