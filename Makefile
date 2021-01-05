.PHONY: build

build:
	sam build

clean:
	rm -f $(wildcard .aws-sam/build/*/*)

run-dev: init-db build dev-sam

dev-sam:
	sam local start-api -p 1946 --env-vars env.json

init-db:
	docker-compose up -d && goose -dir migrations/ postgres "postgres://root:postgres@localhost:5432/postgres?sslmode=disable" up

deploy: clean build
	sam deploy --capabilities CAPABILITY_NAMED_IAM --config-file samconfig.toml --stack-name fender-platform-exercise --s3-bucket "aws-sam-cli-managed-default-samclisourcebucket-1lvgp70bw6awj" --s3-prefix "fender-platform-exercise" --parameter-overrides PostgresURI=$(PG_URL) SigningSecret=$(SUPER_SECRET)