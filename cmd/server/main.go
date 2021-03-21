package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
	"github.com/petarov/appstore-web-search/cmd/common"
)

var (
	ListenAddress string
	ListenPort    int
	WebAppPath    string
)

func init() {
	flag.StringVar(&ListenAddress, "address", "", "Server listen address")
	flag.IntVar(&ListenPort, "port", 5360, "Server listen port")
	flag.StringVar(&WebAppPath, "webapp", "./webapp", "Path to the webapp")
}

func main() {
	fmt.Printf("iTunes Web Search - v%s\n", common.APP_VERSION)
	flag.Parse()

	fmt.Printf("Listening on %s:%d ...\n", ListenAddress, ListenPort)

	WebAppPath, _ = filepath.Abs(WebAppPath)
	if _, err := os.Stat(WebAppPath); err != nil {
		fmt.Printf("Webapp path not found: %v", err)
		return
	}
	fmt.Printf("Serving webapp from: %s ...\n", WebAppPath)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ListenAddress, ListenPort),
		gziphandler.GzipHandler(http.FileServer(http.Dir(WebAppPath))))
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
		return
	}
}
