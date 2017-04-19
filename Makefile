.PHONY: build
build:
	@mkdir -p bin
	@go build -o bin/gloat github.com/gsamokovarov/gloat/gloat

.PHONY: test
test:
	@go test ./...

.PHONY: lint
lint:
	@go vet ./... && golint ./...
