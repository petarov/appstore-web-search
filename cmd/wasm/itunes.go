package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"syscall/js"
	"time"

	"github.com/petarov/itunes-web-search/cmd/common"
)

func search(term string, country string, lang string, media string, entity string, client *http.Client) (status int, json string, err error) {
	query := fmt.Sprintf(
		"https://itunes.apple.com/search?media=%s&entity=%s&term=%s&country=%s&lang=%s&limit=%d&callback=_cb",
		media, entity, term, country, lang, 10)

	fmt.Printf("Search: %s\n", query)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return -1, "", fmt.Errorf("Error creating GET request: %v", err)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s-v%s", common.APP_USER_AGENT, common.APP_VERSION))

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

	// fmt.Printf("STATUS: %d\n", status)
	// fmt.Printf("JSON: %s\n", json)

	search := func() js.Func {
		return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go func() {
				cb := args[1]
				client := &http.Client{Timeout: 5 * time.Second}

				_, json, err := search(args[0].String(), "de", "de_DE", "all", "", client)
				if err != nil {
					fmt.Printf("Error in search: %v", err)
					return
				}
				cb.Invoke(js.Null(), json)
			}()
			return ""
		})
	}
	js.Global().Set("search", search())

	select {}
}
