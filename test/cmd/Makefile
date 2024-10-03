export POSTGRES_DB?=postgres
export POSTGRES_HOST?=localhost
export POSTGRES_PORT?=5432
export POSTGRES_USER?=postgres
export POSTGRES_PASSWORD?=postgres

POSTGRES_DSN?="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

compose-up:
	docker-compose up -d postgres

compose-down:
	docker-compose down

goose-up:
	goose -dir ../../migrations postgres "$(POSTGRES_DSN)" up

goose-down:
	goose -dir ../../migrations postgres "$(POSTGRES_DSN)" down

test-up: compose-up goose-up

test-down: goose-down compose-down