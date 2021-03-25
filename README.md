# App Store Web Search

A web app that queries the Apple App Store in your browser. It uses a Go WebAssembly module to cache results for 60 seconds.

<img src="demo/shot1.png" width="400">

# Installation

`Go` is required. See [installation](https://golang.org/doc/install).

Run `make` to run the server part and open the app in your browser.

Run `make all` to produce a binary package in folder `dist`.

Run `make clean` to clean build artifacts.

# License

[MIT](LICENSE)
