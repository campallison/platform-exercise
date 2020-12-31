.PHONY: build

build:
	sam build

run-dev: build dev-sam run-db

dev-sam:
	sam local start-api -p 1946

run-db:
	docker-compose up -d --no-recreate