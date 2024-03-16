.PHONY: all clean default run build checks pc test

default: checks

run:
	go run cmd/deadnews-template-go/main.go

build:
	go build -o ./dist/ ./...

goreleaser:
	goreleaser --clean --snapshot --skip=publish

install:
	pre-commit install

checks: pc test

pc:
	pre-commit run -a

test:
	go test -v -race -covermode=atomic -coverprofile='coverage.txt' ./...
