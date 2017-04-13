.PHONY: build
build:
	@go build github.com/gsamokovarov/gloat/cmd/gloat

.PHONY: test
test:
	@go test ./...

.PHONY: lint
lint:
	@go vet ./...
