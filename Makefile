GO=go
APP_NAME=pvz-cli-app
BUILD_DIR=build
POSTGRES_CONN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
POSTGRES_TEST_CONN=postgres://postgres:postgres@localhost:5432/postgres_test?sslmode=disable

build: clean
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

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
	goose -dir ./migrations postgres "$(POSTGRES_CONN)" create rename_me sql

goose-up:
	goose -dir ./migrations postgres "$(POSTGRES_CONN)" up

goose-down:
	goose -dir ./migrations postgres "$(POSTGRES_CONN)" down

goose-status:
	goose -dir ./migrations postgres "$(POSTGRES_CONN)" status


# Test build cmds

compose-test-up:
	docker-compose up -d postgres_test

compose-test-down:
	docker-compose down postgres_test

goose-test-up:
	goose -dir ./migrations postgres "$(POSTGRES_TEST_CONN)" up

goose-test-down:
	goose -dir ./migrations postgres "$(POSTGRES_TEST_CONN)" down

test-up: compose-test-up goose-test-up

test-down: compose-test-down goose-test-down