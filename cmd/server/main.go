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
	WebAppPath    string
)

func init() {
	flag.StringVar(&ListenAddress, "address", "localhost", "Server listen address")
	flag.IntVar(&ListenPort, "port", 5360, "Server listen port")
	flag.StringVar(&WebAppPath, "webapp", "./webapp", "Path to the webapp")
}

func main() {
	fmt.Printf("iTunes Web Search - v%s\n", common.APP_VERSION)
	flag.Parse()

	fmt.Printf("Listening on %s:%d ...\n", ListenAddress, ListenPort)

	WebAppPath, _ = filepath.Abs(WebAppPath)
	fmt.Printf("Serving webapp from: %s ...\n", WebAppPath)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ListenAddress, ListenPort),
		http.FileServer(http.Dir(WebAppPath)))
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
		return
	}
}
