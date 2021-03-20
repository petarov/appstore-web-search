package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/petarov/itunes-web-search/cmd/common"
)

var (
	ListenAddress string
	ListenPort    int
	AssetsPath    string
)

func init() {
	flag.StringVar(&ListenAddress, "address", "localhost", "Server listen address")
	flag.IntVar(&ListenPort, "port", 5360, "Server listen port")
	flag.StringVar(&AssetsPath, "assets", "./assets", "Path to web assets to serve")
}

func main() {
	fmt.Println("iTunes Web Search - v", common.APP_VERSION)
	flag.Parse()

	fmt.Printf("Listening on %s:%d ...\n", ListenAddress, ListenPort)

	AssetsPath, _ = filepath.Abs(AssetsPath)
	fmt.Printf("Serving assets from: %s ...\n", AssetsPath)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ListenAddress, ListenPort),
		http.FileServer(http.Dir(AssetsPath)))
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
		return
	}
}
