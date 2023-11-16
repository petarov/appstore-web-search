# App Store Web Search

A [web app](https://vexelon.net/asws) that queries the Apple App Store in your browser. 

Mostly JS + a Go WebAssembly module that caches search results for 60 seconds.

<img src="demo/shot1.png" width="300">

# Installation

`Go` is required. See [installation](https://golang.org/doc/install).

Run `make` to run the server part and open the app in your browser.

Run `make build` to produce the web app files in `webapp`.

Run `make all` to produce a server executable and build files.

Run `make clean` to clean build artifacts.

# License

[MIT](LICENSE)
