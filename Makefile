BINARY_NAME=logsift
GEN_BINARY=loggen

.PHONY: all build test bench clean lint run

all: build

build:
	go build -o $(BINARY_NAME) ./cmd/logsift
	go build -o $(GEN_BINARY) ./cmd/loggen

test:
	go test -race -v ./...

bench:
	go test -benchmem -bench=. ./pkg/...

clean:
	rm -f $(BINARY_NAME) $(GEN_BINARY) coverage.out *.db

lint:
	go vet ./...
