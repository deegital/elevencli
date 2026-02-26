BINARY := elevencli

.PHONY: build install clean vet lint

build:
	go build -o $(BINARY) .

install:
	go install .

clean:
	rm -f $(BINARY)

vet:
	go vet ./...

lint:
	golangci-lint run ./...
