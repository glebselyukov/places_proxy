export GO111MODULE := on
MODULE_NAME:=$(shell sh -c 'cat go.mod | grep module | sed -e "s/module //"')

DOCKER_COMPOSE_PROJECT_FILE=./build/proxy/docker-compose.yml
DOCKER_COMPOSE_CMD=docker-compose -f $(DOCKER_COMPOSE_PROJECT_FILE)

REQUIRED_BINS := goimports go docker docker-compose
$(foreach bin,$(REQUIRED_BINS),\
    $(if $(shell command -v $(bin) 2> /dev/null),\
        $(),\
        $(error Please install "$(bin)", my friend)))

all: prepare format

# PROJECT
project_build:
	$(DOCKER_COMPOSE_CMD) build --force-rm

project_run:
	$(DOCKER_COMPOSE_CMD) up -d

project_down:
	$(DOCKER_COMPOSE_CMD) down

# TESTS
clear:
	rm -f coverage.out
	rm -f coverage.html

tests: clear
	go test -covermode=count -coverprofile=coverage.out `go list ./...` | grep -q ""
	go tool cover -html=coverage.out -o coverage.html

coverage: tests
	go tool cover -func=coverage.out

# FORMAT
format:
	go fmt `go list ./... | grep -v /vendor/`
	goimports -w -local ${MODULE_NAME} `go list -f {{.Dir}} ./...`

# PREPARE
prepare:
	go mod download
	go generate ./...
