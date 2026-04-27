# Load config from .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all build run dev lint migrate seed add-migrate tidy build clean help

# Default target
all: build

## help: Show help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: Run app
run:
	go run $(MAIN_PATH)

## dev: Run app with hot reload
dev:
	air

## lint: Run golangci-lint
lint:
	golangci-lint run

## migrate: Execute all pending migrations
migrate:
	go run $(MIGRATE_PATH)

## seed: Run seeder
seed:
	go run $(SEED_PATH)

## add-migrate: Generate new migration
add-migrate:
	@if [ -z "$(name)" ]; then echo "Error: 'name' is required. Ex: make add-migrate name=CreateUser"; exit 1; fi
	go run $(ADD_MIGRATE_PATH) $(name)

## tidy: Cleaning and verifying go.mod & go.sum
tidy:
	go mod tidy
	go mod verify

## build: Compile app to binary
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

## clean: Delete binary and temp files
clean:
	rm -rf bin/
	rm -rf uploads/temp/*
	rm -rf tmp/
	rm -rf *.log

## dbml: Generate database diagram
dbml:
	PGPASSWORD=$(DB_PASSWORD) pg_dump -s -U $(DB_USER) -h localhost -p $(DB_PORT) $(DB_NAME) > $(DOCS_ERD_PATH)/$(DB_NAME).sql
	sql2dbml $(DOCS_ERD_PATH)/${DB_NAME}.sql --postgres -o $(DOCS_ERD_PATH)/${DB_NAME}.dbml