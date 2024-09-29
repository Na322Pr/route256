POSTGRES_CONN=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable

compose-up:
	docker-compose up -d postgres

compose-down:
	docker-compose down

goose-up:
	goose -dir ../../migrations postgres "$(POSTGRES_CONN)" up

goose-down:
	goose -dir ../../migrations postgres "$(POSTGRES_CONN)" down

test-up: compose-up goose-up

test-down: goose-down compose-down