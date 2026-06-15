BINDIR ?= $(HOME)/.local/bin
BINARY := fogus

.PHONY: build test install clean

build:
	go build -o $(BINARY) .

test:
	go test ./...

install:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(BINARY) .

clean:
	rm -f $(BINARY)
