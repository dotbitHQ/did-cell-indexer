# build file
GO_BUILD=go build -ldflags -s -v

svr: BIN_BINARY_NAME=did_indexer_svr
svr:
	GO111MODULE=on $(GO_BUILD) -o $(BIN_BINARY_NAME) cmd/main.go
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run ./$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"

update:
	go mod tidy
