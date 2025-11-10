BINARY=generate-e2e

BIN_DIR=bin

GO=go

all: build

build:
	$(GO) build -o $(BIN_DIR)/$(BINARY) .

install:
	$(GO) install ./cmd/generate-e2e

run:
	$(GO) run ./cmd/generate-e2e

clean:
	rm -rf $(BIN_DIR)

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(BINARY)-linux ./cmd/generate-e2e

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(BINARY).exe ./cmd/generate-e2e
