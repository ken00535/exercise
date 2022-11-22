GIT_HASH := $(shell git rev-parse HEAD)
GIT_TAG := $(shell git tag -l --sort=-v:refname | head -n 1)
PROJ_NAME := assignment

.PHONY: all test

all:
	source ./.env && go run ./cmd/server

init:
	@cp .env.dev .env
	@cp ./config/app-dev.yml ./config/app.yml

lint:
	golangci-lint run