BIN        = gg
ENTRYPOINT = main.go
LDFLAGS    = "-s -w"
SRCS       = $(shell find . -type f -name '*.go') go.mod go.sum

$(BIN): $(SRCS)
	GO111MODULE=on go build -ldflags $(LDFLAGS) -o $(BIN) $(ENTRYPOINT)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: clean
clean:
	go clean

.PHONY: build
build: $(BIN)
