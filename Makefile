DB_STRING="postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)"
MIGRATIONS_DIR="./migrations"

build:
	@go build -o bin/main cmd/main.go
	@chmod +x bin/main

run:
	@bin/main

all: build up run

up:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up

down:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) down