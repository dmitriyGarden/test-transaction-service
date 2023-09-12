SHELL:=/bin/bash

run:
	go run main.go

migration-build:
	go build -v -o ./bin/migrations ./migrations/

migration-create-sql:
	go build -v -o ./bin/migrations ./migrations/ && \
    ./bin/migrations create $(name) sql

migration-create-go:
	go build -v -o ./bin/migrations ./migrations/ && \
    ./bin/migrations create $(name) go

migration-status:
	go build -v -o ./bin/migrations ./migrations/ && \
    ./bin/migrations status

migration-up:
	go build -v -o ./bin/migrations ./migrations/ && \
    ./bin/migrations up

migration-up-by-one:
	go build -v -o ./bin/migrations ./migrations/ && \
    ./bin/migrations up-by-one

migration-down-to:
	go build -v -o ./bin/migrations ./migrations/ && \
    ./bin/migrations down-to $(version)
