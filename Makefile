#!make

SWAGGER_YAML:=internal/docs/swagger.yaml
GO_BINARY:=go

.PHONY: test

setup:
	./scripts/setup.sh

clean: 
	rm -f ./bin/*

build: api-server-build-mac-arm64

test:
	 $(GO_BINARY) test --v ./...

test-coverage:
	 $(GO_BINARY) test -cover --v ./...

test-int:
	$(GO_BINARY) test -v ./... -p 1 -tags=integration

# Api Server commands
api-server-build: api-server-build-mac-arm64

api-server-build-mac-arm64:
	GOOS=darwin GOARCH=arm64 $(GO_BINARY) build -o bin/yuka-api-server-mac-arm64 ./cmd/apiserver

api-server-build-mac-amd64:
	GOOS=darwin GOARCH=amd64 $(GO_BINARY) build -o bin/yuka-api-server-mac-amd64 ./cmd/apiserver

api-server-build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO_BINARY) build -o bin/yuka-api-server-linux-amd64 ./cmd/apiserver

api-server-build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO_BINARY) build -o bin/yuka-api-server-linux-arm64 ./cmd/apiserver

api-server-build-windows:
	GOOS=windows GOARCH=amd64 $(GO_BINARY) build -o bin/yuka-api-server-windows.exe ./cmd/apiserver
# api-server-build-all: build-mac-arm64 build-mac-amd64 build-linux-amd64 build-linux-arm64 build-windows

api-server-start-dev: gen-docs api-server-build
	ENVIRONMENT=local LOG_LEVEL=debug DATABASE_HOSTNAME="localhost" DATABASE_USERNAME="postgres" DATABASE_PASSWORD="password" DATABASE_NAME="mydb" DATABASE_PORT=5432 APISERVER_ADDRESS="localhost" ./bin/yuka-api-server-mac-arm64

api-server-start: gen-docs api-server-build
	LOG_LEVEL=debug DATABASE_HOSTNAME="localhost" DATABASE_USERNAME="postgres" DATABASE_PASSWORD="password" DATABASE_NAME="mydb" DATABASE_PORT=5432 APISERVER_ADDRESS="localhost" ./bin/yuka-api-server-mac-arm64

# yukactl commands
yukactl-build: yukactl-build-mac-arm64

yukactl-build-mac-arm64:
	GOOS=darwin GOARCH=arm64 $(GO_BINARY) build -o bin/yukactl-mac-arm64 ./cmd/yukactl/

yukactl-build-mac-amd64:
	GOOS=darwin GOARCH=amd64 $(GO_BINARY) build -o bin/yukactl-mac-amd64 ./cmd/yukactl/

yukactl-build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO_BINARY) build -o bin/yukactl-linux-amd64 ./cmd/yukactl/

yukactl-build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO_BINARY) build -o bin/yukactl-linux-arm64 ./cmd/yukactl/

yukactl-build-windows:
	GOOS=windows GOARCH=amd64 $(GO_BINARY) build -o bin/yukactl-windows.exe ./cmd/yukactl/


# Docs commands	
gen-docs:
	$$HOME/go/bin/swag init -g ./cmd/apiserver/main.go -o ./docs

gen-client: gen-docs
	swagger generate client --target internal/api --spec ./docs/swagger.json --name yuka --client-package api-clients --model-package api-models
