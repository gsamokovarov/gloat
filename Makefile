.PHONY: build
build:
	@mkdir -p bin
	@go build -o bin/gloat github.com/gsamokovarov/gloat/gloat

.PHONY: test
test:
	@go test ./...

.PHONY: test
test.sqlite:
	@env DATABASE_URL=sqlite3://:memory: go test ./...

.PHONY: lint
lint:
	@go vet ./... && golint ./...
