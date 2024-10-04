GO=go
APP_NAME=pvz-cli-app
BUILD_DIR=build

export POSTGRES_DB?=postgres
export POSTGRES_HOST?=localhost
export POSTGRES_PORT?=5432
export POSTGRES_USER?=postgres
export POSTGRES_PASSWORD?=postgres

POSTGRES_DSN?=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

CONFIG_PATH=./config/config.yaml

build: clean
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME) --config="$(CONFIG_PATH)"

install:
	@echo "Installing dependencies..."
	$(GO) mod tidy
	$(GO) install github.com/uudashr/gocognit/cmd/gocognit@latest
	$(GO) install github.com/fzipp/gocyclo/cmd/gocyclo@latest

update:
	@echo "Updating dependencies..."
	$(GO) get -u ./...

gocognit:
	@echo "Checking gocognit..."
	gocognit -over 5 .

gocyclo:
	@echo "Checking gocyclo..."
	gocyclo -over 5 .

test:
	$(GO) test ./...

coverage: 
	$(GO) test -coverprofile=coverage.out ./...

coverage_html: coverage
	$(GO) tool cover -html=coverage.out

clean: 
	rm -rf $(BUILD_DIR)/$(APP_NAME)

compose-up:
	docker-compose up -d postgres

compose-down:
	docker-compose down

compose-stop:
	docker-compose stop postgres

compose-start:
	docker-compose start postgres

compose-ps:
	docker-compose ps postgres

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" create rename_me sql

goose-up:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" up

goose-down:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" down

goose-status:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" status
