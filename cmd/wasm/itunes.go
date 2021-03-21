package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"
	"time"

	"github.com/petarov/itunes-web-search/cmd/common"
)

func search(term string, country string, lang string, media string, entity string, client *http.Client) (json string, err error) {
	query := fmt.Sprintf(
		"https://itunes.apple.com/search?media=%s&term=%s&country=%s&limit=%d&callback=_cb",
		media, term, country, 20)

	if lang != "" {
		query = fmt.Sprintf("%s&lang=%s", query, lang)
	}
	if entity != "" {
		query = fmt.Sprintf("%s&entity=%s", query, entity)
	}

	fmt.Printf("Search: %s\n", query)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating GET request: %v", err)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s-v%s", common.APP_USER_AGENT, common.APP_VERSION))

	// reqOut, err := httputil.DumpRequest(req, false)
	// fmt.Println(string(reqOut))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error in HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading HTTP response body: %v", err)
	}

	if resp.StatusCode/100 != 2 {
		return string(body), fmt.Errorf("HTTP Error %d", resp.StatusCode)
	}

	return string(body), nil
}

func main() {
	fmt.Println("*** Welcome to App Store Web Search ***")

	search := func() js.Func {
		return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go func() {
				term := args[0].String()
				country := args[1].String()
				media := args[2].String()
				cb := args[3]

				json, err := search(term, country, "", media, "",
					&http.Client{Timeout: 4 * time.Second})
				if err != nil {
					cb.Invoke(err.Error(), json)
					return
				} else {
					cb.Invoke(js.Null(), json)
				}
			}()
			return nil
		})
	}
	js.Global().Set("search", search())

	select {}
}
