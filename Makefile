BINARY := prometheus
LDFLAGS := -ldflags="-s -w"
CGO := CGO_ENABLED=0

.PHONY: build test fmt lint clean size-check

build:
	$(CGO) go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/prometheus

test:
	$(CGO) go test -race -count=1 -timeout=120s ./...

fmt:
	gofmt -w ./cmd ./internal

lint:
	@echo "golangci-lint configuration scaffolded in .golangci.yml"

size-check:
	@echo "Size check placeholder: requires embedded llama-server artifacts"

clean:
	@if [ -d bin ]; then rm -rf bin; fi

