ifeq ($(OS),Windows_NT)
	BROWSER = start
	OS=windows
else
	UNAME := $(shell uname)
	ifeq ($(UNAME), Darwin)
		BROWSER = open
		OS=darwin
	else
		BROWSER = xdg-open
		OS=linux
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

.PHONY: all

serve: cp_wasm appstore.wasm open
build: cp_wasm appstore.wasm server.cmd
all: clean build dist

cp_wasm:
	test -f webapp/wasm_exec.js || cp $(GO_ROOT)/misc/wasm/wasm_exec.js webapp/

%.wasm: cmd/wasm/%.go
	GOOS=js GOARCH=wasm go generate
	GOOS=js GOARCH=wasm go build -o "$@" "$<"
	mv appstore.wasm $(PROJECT_DIR)/webapp/

%.cmd: cmd/%/main.go
	GOOS=$(OS) GOARCH=amd64 go build -o asws_server_$(OS)_amd64 "$<"

open:
	$(BROWSER) 'http://$(SERVE_ADDR):$(SERVE_PORT)'
	go run cmd/server/main.go -port $(SERVE_PORT) -address $(SERVE_ADDR)

dist:
	test -d dist || mkdir -p dist/webapp
	cp webapp/* dist/webapp
	mv asws_server_* dist/

clean:
	rm -f webapp/*.wasm
	rm -f asws_server_*
	test -d dist && rm -f dist/webapp/* && rm -f dist/asws_server_* && rmdir -p dist/webapp