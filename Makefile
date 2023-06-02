.PHONY: build run clean

GOROOT = $(shell go env GOROOT)

build:
	GOOS=js GOARCH=wasm go build -o main.wasm
	cp $(GOROOT)/misc/wasm/wasm_exec.js .

run:
	python3 -m http.server 8080

clean:
	rm -f main.wasm wasm_exec.js
