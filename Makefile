BUILD_PATH=tmp
GEN_PATH=internal/gen
PROTO_PATH=proto
MAIN_PACKAGE_PATH=./cmd/protoc-gen-sqlc
BINARY_NAME=protoc-gen-sqlc

LIBRARY_EXAMPLE_PATH=examples/library
LIBRARY_EXAMPLE_CONFIG=buf.gen.yaml

BUF_VERSION=v1.47.2
SQLC_VERSION=v1.27.0
GOLANGCILINT_VERSION=v2.0.2
PROTOC_GEN_GO_VERSION=v1.35.2

.PHONY: help
help:
	@echo 'Build scripts for $(BINARY_NAME)'
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

.PHONY: format
format:
	go fmt ./...
	go run mvdan.cc/gofumpt@latest -w cmd internal

.PHONY: tidy
tidy: format
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit: tidy
	go mod verify
	go vet ./...
	go test -v -race -failfast -buildvcs -vet=off ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCILINT_VERSION) run ./...
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) lint

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=$(BUILD_PATH)/coverage.out ./...
	go tool cover -html=$(BUILD_PATH)/coverage.out

## build: build the executable
.PHONY: build
build: protoc-gen-go
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) generate --path $(PROTO_PATH)/sqlc
	go build -o=$(BUILD_PATH)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

## run: run the executable
.PHONY: run
run: build
	$(BUILD_PATH)/$(BINARY_NAME)

## protoc: compile protocol buffers
.PHONY: protoc
protoc: build
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) generate \
		--template $(LIBRARY_EXAMPLE_PATH)/$(LIBRARY_EXAMPLE_CONFIG) \
		--path $(PROTO_PATH)/examples

protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

## sqlc: generate idiomatic code from SQL files
.PHONY: sqlc
sqlc: protoc
	cd $(LIBRARY_EXAMPLE_PATH) && \
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate

## exampledb/up: start PostgreSQL container used for the example
.PHONY: exampledb/up
exampledb/up:
	docker compose -f $(LIBRARY_EXAMPLE_PATH)/compose.yml up -d

## exampledb/down: stop PostgreSQL container used for the example
.PHONY: exampledb/down
exampledb/down:
	docker compose -f $(LIBRARY_EXAMPLE_PATH)/compose.yml down

## example/library: run library example
.PHONY: example/library
example/library: exampledb/up sqlc
	cd $(LIBRARY_EXAMPLE_PATH) && go run main.go

## push: push changes to the remote Git repository
.PHONY: push
push: clean sqlc audit test/cover no-dirty
	git push

## upgrade: upgrade Go dependencies
.PHONY: upgrade
upgrade:
	go get -u -t ./...
	go mod tidy -v

## clean: clean all generated artifacts
.PHONY: clean
clean: exampledb/down
	rm -rf $(BUILD_PATH)
	rm -rf $(LIBRARY_EXAMPLE_PATH)/$(GEN_PATH)
