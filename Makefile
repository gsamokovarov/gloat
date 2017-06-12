.PHONY: build
build:
	@mkdir -p bin
	@go build -o bin/gloat github.com/gsamokovarov/gloat/cmd/gloat

.PHONY: test
test: embed
	@go test ./...

.PHONY: test.sqlite
test.sqlite: embed
	@env DATABASE_URL=sqlite3://:memory: go test ./...

.PHONY: embed
embed:
	@go-bindata -pkg gloat -o test_assets.go testdata/migrations/*

.PHONY: lint
lint:
	@go vet ./... && golint ./...
