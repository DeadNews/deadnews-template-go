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

bumped:
	git cliff --bumped-version

# make release-tag_name
# make release-$(git cliff --bumped-version)-alpha.0
release-%: checks
	git cliff -o CHANGELOG.md --tag $*
	pre-commit run --files CHANGELOG.md || pre-commit run --files CHANGELOG.md
	git add CHANGELOG.md
	git commit -m "chore(release): prepare for $*"
	git push
	git tag -a $* -m "chore(release): $*"
	git push origin $*
	git tag --verify $*
