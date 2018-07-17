.PHONY: build  package start-docker
DIST_PATH = dist
RELEASE_PATH = $(DIST_PATH)/release
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse HEAD)
GO = go
GINKGO = ginkgo
GO_LINKER_FLAGS ?= -ldflags \
	   "-X 'luckgo/model.BuildDate=$(BUILD_DATE)' \
	   -X luckgo/model.BuildHash=$(BUILD_HASH)"
BUILDER_GOOS_GOARCH=$(shell $(GO) env GOOS)_$(shell $(GO) env GOARCH)
PACKAGESLISTS=$(shell $(GO) list ./...)
TESTFLAGS ?= -short
PACKAGESLISTS_COMMA=$(shell echo $(PACKAGESLISTS) | tr ' ' ',')
ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(GO) install -i $(GOFLAGS) $(GO_LINKER_FLAGS) ./...

package:
	@echo Packaging luckgo
	rm -rf $(DIST_PATH)
	mkdir -p $(RELEASE_PATH)/bin
	mkdir -p $(RELEASE_PATH)/config
	mkdir -p $(RELEASE_PATH)/logs
	cp $(GOBIN)/luckgo $(RELEASE_PATH)/bin
	cp config/*.json $(RELEASE_PATH)/config
	tar -C $(DIST_PATH) -czf $(RELEASE_PATH)-$(BUILDER_GOOS_GOARCH).tar.gz release

govet: ## Runs govet against all packages.
	@echo Running GOVET
	$(GO) vet $(GOFLAGS) $(PACKAGESLISTS) || exit 1

gofmt: ## Runs gofmt against all packages.
	@echo Running GOFMT
	@echo $(PACKAGESLISTS)
	@for package in $(PACKAGESLISTS);do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success"; \

test: clean-docker start-docker test-server

test-server:
	docker run --rm  --name ginkgo-test --network my-net -v $(GOPATH)/src/luckgo:/go/src/luckgo  golang:latest /bin/sh -c  "/go/src/luckgo/ginkgo -r -trace -cover  -coverprofile=coverprofile.txt -outputdir=/go/src/luckgo   /go/src/luckgo"

start-docker: ## Starts the docker containers for local development.
	@echo Starting docker containers

	@if [ $(shell docker ps -a | grep -ci critic-mysql) -eq 0 ]; then \
		echo starting luckgo-mysql; \
		docker run --network my-net --name luckgo-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 \
		-e MYSQL_USER=luckgo -e MYSQL_PASSWORD=123456 -e MYSQL_DATABASE=luckgo_test -d mysql:5.7 > /dev/null; \
	elif [ $(shell docker ps | grep -ci critic-mysql) -eq 0 ]; then \
		echo restarting critic-mysql; \
		docker start critic-mysql > /dev/null; \
	fi

stop-docker: ## Stops the docker containers for local development.
	@echo Stopping docker containers

	@if [ $(shell docker ps -a | grep -ci luckgo-mysql) -eq 1 ]; then \
		echo stopping critic-mysql; \
		docker stop critic-mysql > /dev/null; \
	fi

clean-docker: ## Deletes the docker containers for local development.
	@echo Removing docker containers

	@if [ $(shell docker ps -a | grep -ci luckgo-mysql) -eq 1 ]; then \
		echo removing critic-mysql; \
		docker stop critic-mysql > /dev/null; \
		docker rm -v critic-mysql > /dev/null; \
	fi

