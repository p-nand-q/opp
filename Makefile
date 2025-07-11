.PHONY: all build test clean install

all: test build

build:
	go build -o bin/opp ./cmd/opp

test:
	go test -v ./...

clean:
	rm -rf bin/

install:
	go install ./cmd/opp

# Example usage
example-debug:
	go run ./cmd/opp -D DEBUG examples/hello.opp.go

example-windows:
	go run ./cmd/opp -D WINDOWS examples/hello.opp.go