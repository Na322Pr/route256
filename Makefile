GO=go
APP_NAME=pvz-cli-app
BUILD_DIR=$(CURDIR)/build

OUT_PATH:=$(CURDIR)/pkg
LOCAL_BIN:=$(CURDIR)/bin

CONFIG_PATH=./config/config.yaml


# -----------------------
# 
# -----------------------

POSTGRES_DB=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres

POSTGRES_DSN=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable


# -----------------------
# Запуск приложения на Go
# -----------------------

build: clean
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go

run:
# ./$(BUILD_DIR)/$(APP_NAME) --config="$(CONFIG_PATH)"
	go run ./cmd/pvz-service/main.go --config="$(CONFIG_PATH)"

run-cli: 
	go run ./cmd/pvz-cli/main.go --config="$(CONFIG_PATH)"

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

# ---------------------------
# Запуск базы данных в Docker
# ---------------------------

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

# ---------------------------
# Запуск миграций через Goose
# ---------------------------

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


# ---------------------------------
# Запуск кодогенерации через protoc
# ---------------------------------

bin-deps: .vendor-proto
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest

generate:
	mkdir -p ${OUT_PATH}
	protoc --proto_path api --proto_path vendor.protogen \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=${OUT_PATH} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=${OUT_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ${OUT_PATH} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=${OUT_PATH} \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:${OUT_PATH}" \
		./api/pvz-service/v1/pvz_service.proto

.vendor-proto: .vendor-proto/google/protobuf .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options .vendor-proto/validate

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

.vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
		cd vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
		rm -rf vendor.protogen/protobuf

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p  vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.protogen/validate
		mv vendor.protogen/tmp/validate vendor.protogen/
		rm -rf vendor.protogen/tmp

