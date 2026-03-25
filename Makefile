include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	docker compose up -d postgres-db

env-down:
	docker compose down postgres-db

env-recover:
	docker compose down postgres-db && rm -rf out/pgdata

migrate-create:
	docker compose run --rm migrator create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"

migrate-up:
	docker compose run --rm migrator \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres-db:5432/${POSTGRES_DB}?sslmode=disable \
	up

migrate-down:
	docker compose run --rm migrator \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres-db:5432/${POSTGRES_DB}?sslmode=disable \
	down

