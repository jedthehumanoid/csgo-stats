.PHONY: all
all: build

.PHONY: build
build:
	go build
	cd svelte && npm run check 
	cd svelte && npm run build

.PHONY: fmt
fmt:
	gofmt -l -s -w .
	cd svelte && npm run fmt

