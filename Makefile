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

BUILD_DIR   = $(shell pwd)
GO_ROOT     = $(shell go env GOROOT)
SERVE_PORT  = 5360

.PHONY: all clean serve

build: cp_wasm itunes.wasm
all: build serve

cp_wasm:
	test -f webapp/wasm_exec.js || cp $(GO_ROOT)/misc/wasm/wasm_exec.js webapp/

%.wasm: cmd/wasm/%.go
	GOOS=js GOARCH=wasm go generate
	GOOS=js GOARCH=wasm go build -o "$@" "$<"
	mv itunes.wasm $(BUILD_DIR)/webapp/

serve:
	$(BROWSER) 'http://localhost:$(SERVE_PORT)'
	go run cmd/server/main.go -port $(SERVE_PORT) -address local.petrovs.net

clean:
	rm -f webapp/*.wasm