.PHONY: build

build:
	sam build

run-dev: init-db build dev-sam

dev-sam:
	sam local start-api -p 1946 --env-vars env.json

init-db:
	docker-compose up -d && goose -dir migrations/ postgres "postgres://root:postgres@localhost:5432/postgres?sslmode=disable" up
