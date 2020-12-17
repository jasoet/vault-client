.EXPORT_ALL_VARIABLES:
OUT_DIR := ./_output
BIN_DIR := ./bin

APP_NAME=vault-client
PACKAGE=github.com/jasoet/vault-client
CURRENT_DIR=$(shell pwd)

DOCKER_COMPOSE_FILE=$(CURRENT_DIR)/test/docker-compose.yaml

ifneq (${GITHUB_WORKFLOW},)
TEST_VAULT_ADDR=http://vault:8200
else
TEST_VAULT_ADDR=http://127.0.0.1:18200
endif

TEST_VAULT_TOKEN=localhost

VERSION=$(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell if [ -z "`git status --porcelain`" ]; then git describe --exact-match --tags HEAD 2>/dev/null; fi)

$(shell mkdir -p $(OUT_DIR) $(BIN_DIR))

# perform static compilation
STATIC_BUILD?=true

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT}

ifeq (${STATIC_BUILD}, true)
override LDFLAGS += -extldflags "-static"
endif

ifneq (${GIT_TAG},)
IMAGE_TAG=${GIT_TAG}
#IMAGE_TRACK=stable
LDFLAGS += -X ${PACKAGE}.gitTag=${GIT_TAG}
else
IMAGE_TAG?=$(GIT_COMMIT)
#IMAGE_TRACK=latest
endif

# Code build targets
.PHONY: vendor
vendor:
	go mod vendor

#.PHONY: build.binaries
#build.binaries:
#	CGO_ENABLED=0 GO111MODULE=on go build -mod=vendor -a -ldflags '${LDFLAGS}' -o ${BIN_DIR}/gopay-cd ./cmd/gopay-cd/main.go

#.PHONY: build
#build: vendor build.binaries

# Docker Compose Integration Test tasks
.PHONY:  compose-up
compose-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY:  compose-destroy
compose-destroy:
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop
	docker-compose -f $(DOCKER_COMPOSE_FILE) rm -f

.PHONY:  compose-recreate
compose-recreate: compose-destroy compose-up

# Main Test Targets (without docker)
.PHONY: test
test:
	go test -race -coverprofile=$(OUT_DIR)/coverage.out ./...

.PHONY: integration-test
integration-test:
	TEST_VAULT_ADDR=$(TEST_VAULT_ADDR) TEST_VAULT_TOKEN=$(TEST_VAULT_TOKEN) go test -race -tags=integration -coverprofile=$(OUT_DIR)/coverage.out ./...

.PHONY: docker-integration-test
docker-integration-test: compose-up integration-test
