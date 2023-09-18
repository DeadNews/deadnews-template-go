run_app:
	go run cmd/deadnews-template-go/main.go

test:
	go test -v -race -covermode=atomic -coverprofile='coverage.txt' ./...

build:
	go build -o ./dist/ ./...

.PHONY: test
