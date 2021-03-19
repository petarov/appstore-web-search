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

.PHONY: all clean serve

all: main.wasm serve

%.wasm: %.go
	GOOS=js GOARCH=wasm go generate
	GOOS=js GOARCH=wasm go build -o "$@" "$<"

serve:
	$(BROWSER) 'http://localhost:5360'
	python3 -m http.server 5360

clean:
	rm -f *.wasm