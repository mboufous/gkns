# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /usr/bin/env bash

WORKDIR := $(shell pwd)
APP_NAME := gkns

.PHONY: build
# build executable file for dev
build:
	@go build -o $(WORKDIR)/output/bin/$(APP_NAME)

.PHONY: run
# run executable file
run:
	@$(WORKDIR)/output/bin/$(APP_NAME)

.PHONY: clean
# clean build cache and docker images
clean:
	 rm -rf output

.PHONY: release
## release a version and push to github
#release:
#	goreleaser release --rm-dist

# show help
help:
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\w0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.PHONY: all
all: clean build run

.DEFAULT_GOAL := all