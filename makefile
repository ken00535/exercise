GIT_HASH := $(shell git rev-parse HEAD)
GIT_TAG := $(shell git tag -l --sort=-v:refname | head -n 1)
PROJ_NAME := shorten

.PHONY: all test

all:
	source ./.env && go run ./cmd/server

init:
	@cp .env.dev .env
	@cp ./config/app-dev.yml ./config/app.yml

lint:
	golangci-lint run

test.fuzz:
	go test -fuzz=Fuzz -fuzztime 20s ./internal/app/usecase

test.bench:
	go test -v -bench=. -benchtime=1000x ./test/benchmark

docker.up:
	docker-compose up -d

docker.down:
	docker-compose down

docker.start:
	docker-compose start -d

docker.stop:
	docker.compose stop

db.up:
	source ./.goose.sh && goose -dir deployments/database up

db.down:
	source ./.goose.sh && goose -dir deployments/database down

db.mock:
	source ./.pgenv.sh && psql -f ./test/testdata/mock.sql

mock.update:
	go generate ./...