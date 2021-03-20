ifeq ($(OS),Windows_NT)
	BROWSER = start
else
	UNAME := $(shell uname)
	ifeq ($(UNAME), Linux)
		BROWSER = xdg-open
	else
		BROWSER = open
	endif
endif

BUILD_DIR   = $(shell pwd)
GO_ROOT     = $(shell go env GOROOT)
SERVE_PORT  = 80

.PHONY: all clean serve

build: cp_wasm itunes.wasm
all: build serve

cp_wasm:
	test -f assets/wasm_exec.js || cp $(GO_ROOT)/misc/wasm/wasm_exec.js assets/

%.wasm: cmd/wasm/%.go
	GOOS=js GOARCH=wasm go generate
	GOOS=js GOARCH=wasm go build -o "$@" "$<"
	mv itunes.wasm $(BUILD_DIR)/assets/

serve:
	$(BROWSER) 'http://localhost:$(SERVE_PORT)'
	go run cmd/server/main.go -port $(SERVE_PORT)

clean:
	rm -f assets/*.wasm