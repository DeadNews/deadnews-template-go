.PHONY: build test run

default: build

run:
	go run cmd/deadnews-template-go/main.go

build:
	go build -o ./dist/ ./...

pc-install:
	pre-commit install

checks: pc-run test

test:
	go test -v -race -covermode=atomic -coverprofile='coverage.txt' ./...

pc-run:
	pre-commit run -a

docker: compose-up

compose-up:
	docker compose up --build

compose-down:
	docker compose down
