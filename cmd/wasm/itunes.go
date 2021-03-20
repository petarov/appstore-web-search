package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"syscall/js"
	"time"

	"golang.org/x/net/http2"

	"github.com/petarov/itunes-web-search/cmd/common"
)

func search(term string, country string, lang string, media string, entity string, client *http.Client) (status int, json string, err error) {
	query := fmt.Sprintf(
		"https://itunes.apple.com/search?media=%s&entity=%s&term=%s&country=%s&lang=%s&limit=%d",
		media, entity, term, country, lang, 25)

	fmt.Printf("Search: %s\n", query)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return -1, "", fmt.Errorf("Error creating GET request: %v", err)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("itunes-web-search-v%s", common.APP_VERSION))

	reqOut, err := httputil.DumpRequest(req, false)
	fmt.Println(string(reqOut))

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", fmt.Errorf("Error in HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", fmt.Errorf("Error reading HTTP response body: %v", err)
	}

	return resp.StatusCode, string(body), nil
}

func main() {
	fmt.Println("*** Welcome to iTunes Web Search ***")

	// status, json, err := search("midpoints", "de", "de_DE", "all", "", client)
	// if err != nil {
	// 	fmt.Printf("Error in search: %v", err)
	// 	return
	// }

	// fmt.Printf("STATUS: %d\n", status)
	// fmt.Printf("JSON: %s\n", json)

	// stopButton := js.Global().Get("document").Call("getElementById", "stop")
	// stopButton.Set("onclick", js.FuncOf(func(js.Value, []js.Value) interface{} {
	// 	go func() {
	// 		client := &http.Client{Timeout: 5 * time.Second}
	// 		client.Transport = &http2.Transport{TLSClientConfig: &tls.Config{}}

	// 		_, json, err := search("midpoints", "de", "de_DE", "all", "", client)
	// 		if err != nil {
	// 			fmt.Printf("Error in search: %v", err)
	// 			return
	// 		}

	// 		fmt.Printf("JSON: %s\n", json)
	// 	}()

	// 	return nil
	// }))

	example := func() js.Func {
		return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go func() {
				client := &http.Client{Timeout: 5 * time.Second}
				client.Transport = &http2.Transport{TLSClientConfig: &tls.Config{}}

				_, json, err := search("midpoints", "de", "de_DE", "all", "", client)
				if err != nil {
					fmt.Printf("Error in search: %v", err)
					return
				}

				fmt.Printf("JSON: %s\n", json)
			}()
			return ""
		})
	}
	js.Global().Set("example", example())

	select {}
}
