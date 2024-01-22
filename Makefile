.PHONY: all clean test default checks docker

default: build

run:
	go run cmd/deadnews-template-go/main.go

build:
	go build -o ./dist/ ./...

pc-install:
	pre-commit install

checks: pc-run test

pc-run:
	pre-commit run -a

test:
	go test -v -race -covermode=atomic -coverprofile='coverage.txt' ./...

docker: compose-up

compose-up:
	docker compose up --build

compose-down:
	docker compose down

goreleaser:
	goreleaser --clean --snapshot --skip=publish

# confirm
semantic-release:
	git tag "$(svu next)"
	git push origin "$(svu next)"

semantic-prerelease:
	git tag "$(svu prerelease --pre-release beta)"
	git push origin "$(svu prerelease --pre-release beta)"

# svu prerelease --pre-release beta.0
