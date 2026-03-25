include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	docker compose up -d postgres-db

env-down:
	docker compose down postgres-db

env-clear:
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

app-up:
	docker compose up -d --build app

app-down:
	docker compose exec -T app kill -TERM 1 &&\
	docker compose down app