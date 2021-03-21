ifeq ($(OS),Windows_NT)
	BROWSER = start
else
	UNAME := $(shell uname)
	ifeq ($(UNAME), Darwin)
		BROWSER = open
	else
		BROWSER = xdg-open
	endif
endif

PROJECT_DIR	= $(shell pwd)
GO_ROOT     = $(shell go env GOROOT)
SERVE_PORT  = 5360
ifeq ($(addr),)
	SERVE_ADDR = localhost
else
	SERVE_ADDR = $(addr)
endif

.PHONY: all clean serve

build: cp_wasm appstore.wasm
all: build serve

cp_wasm:
	test -f webapp/wasm_exec.js || cp $(GO_ROOT)/misc/wasm/wasm_exec.js webapp/

%.wasm: cmd/wasm/%.go
	GOOS=js GOARCH=wasm go generate
	GOOS=js GOARCH=wasm go build -o "$@" "$<"
	mv appstore.wasm $(PROJECT_DIR)/webapp/

serve:
	$(BROWSER) 'http://$(SERVE_ADDR):$(SERVE_PORT)'
	go run cmd/server/main.go -port $(SERVE_PORT) -address $(SERVE_ADDR)

clean:
	rm -f webapp/*.wasm