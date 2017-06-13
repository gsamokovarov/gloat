.PHONY: build
build:
	@mkdir -p bin
	@go build -o bin/gloat github.com/gsamokovarov/gloat/cmd/gloat

.PHONY: test
test:
	@go test ./...

.PHONY: test.sqlite
test.sqlite:
	@env DATABASE_URL=sqlite3://:memory: go test ./...

.PHONY: test.assets
test.assets:
	@go-bindata -pkg gloat -o assets_test.go testdata/migrations/*

.PHONY: lint
lint:
	@go vet ./... && golint ./...
