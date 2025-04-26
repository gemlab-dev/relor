GOPATH := $(shell go env GOPATH)/bin

export PATH := $(PATH):$(GOPATH)

BIN_DIR = bin
PB_DIR = gen/pb

all: clean proto tidy test build 

build:
	go build -o $(BIN_DIR)/ -v ./...

cov:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out

test:
	go test -v ./...

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -rf $(PB_DIR)/*
	#TODO remove all docker containers

tidy:
	go mod tidy
	go mod vendor
	go vet ./...

proto:
	docker build -f Dockerfile.protoc -t protoc-tool .
	docker run  --rm -v .:/src -w /src protoc-tool \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		graph/*.proto
	docker run  --rm -v .:/src -w /src protoc-tool \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		api/*.proto
	docker run  --rm -v .:/src -w /src protoc-tool \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		config/*.proto
	docker run  --rm -v .:/src -w /src protoc-tool \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		db/*.proto

.PHONY: all build test clean cov tidy proto