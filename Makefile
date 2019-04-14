export GO111MODULE := on
MODULE_NAME:=$(shell sh -c 'cat go.mod | grep module | sed -e "s/module //"')

PROXY_SERVICE_ENV=./build/proxy/config/proxy.env

DOCKER_COMPOSE_PROJECT_FILE=./build/proxy/docker-compose.yml
DOCKER_COMPOSE_CMD=docker-compose -f $(DOCKER_COMPOSE_PROJECT_FILE)

REQUIRED_BINS := goimports go docker docker-compose
$(foreach bin,$(REQUIRED_BINS),\
    $(if $(shell command -v $(bin) 2> /dev/null),\
        $(),\
        $(error Please install "$(bin)", my friend)))

.PHONY: all
all: prepare

# PROJECT
.PHONY: project_build
project_build:
	$(DOCKER_COMPOSE_CMD) build --force-rm

.PHONY: project_run
project_run:
	$(DOCKER_COMPOSE_CMD) up -d --force-recreate

.PHONY: project_down
project_down:
	$(DOCKER_COMPOSE_CMD) down

# TESTS
.PHONY: clear
clear:
	rm -f coverage.out
	rm -f coverage.html

.PHONY: tests
tests: clear
	go test -covermode=count -coverprofile=coverage.out `go list ./...` | grep -q ""
	go tool cover -html=coverage.out -o coverage.html

.PHONY: coverage
coverage: tests
	go tool cover -func=coverage.out

# FORMAT
.PHONY: format
format:
	go fmt `go list ./... | grep -v /vendor/`
	goimports -w -local ${MODULE_NAME} `go list -f {{.Dir}} ./...`

# PREPARE
.PHONY: prepare
prepare:
	go mod download
	go generate ./...
