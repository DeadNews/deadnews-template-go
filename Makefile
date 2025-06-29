.PHONY: all clean default run build update checks pc test integr

default: checks

run:
	go run ./cmd/template-go

build:
	go build -o ./dist/ ./...

goreleaser:
	goreleaser --clean --snapshot --skip=publish

update:
	go get -u ./...
	go mod tidy
	go mod verify

checks: pc test
pc:
	pre-commit run -a
test:
	go test -v -race -covermode=atomic -coverprofile='coverage.txt' ./...

integr:
	TESTCONTAINERS=1 go test -v -race -covermode=atomic -coverprofile=coverage.txt ./...

bumped:
	git cliff --bumped-version

# make release TAG=$(git cliff --bumped-version)-alpha.0
release: checks
	git cliff -o CHANGELOG.md --tag $(TAG)
	pre-commit run --files CHANGELOG.md || pre-commit run --files CHANGELOG.md
	git add CHANGELOG.md
	git commit -m "chore(release): prepare for $(TAG)"
	git push
	git tag -a $(TAG) -m "chore(release): $(TAG)"
	git push origin $(TAG)
