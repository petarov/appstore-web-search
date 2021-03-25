package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"syscall/js"
	"time"

	"github.com/petarov/appstore-web-search/cmd/common"
)

type cacheEntry struct {
	data      string
	expiresOn time.Time
}

const (
	CACHE_PURGE_THRESHOLD = 5  // purge after this many entries
	CACHE_TTL             = 60 // seconds
)

var (
	cache                   map[string]*cacheEntry
	concurrentPurgeRoutines chan struct{}
)

// func search(term string, country string, lang string, media string, entity string, client *http.Client) (json string, err error) {
// 	query := fmt.Sprintf(
// 		"https://itunes.apple.com/search?media=%s&term=%s&country=%s&limit=%d&callback=_cb",
// 		media, term, country, 20)

// 	if lang != "" {
// 		query = fmt.Sprintf("%s&lang=%s", query, lang)
// 	}
// 	if entity != "" {
// 		query = fmt.Sprintf("%s&entity=%s", query, entity)
// 	}

// 	fmt.Printf("Search: %s\n", query)

// 	req, err := http.NewRequest("GET", query, nil)
// 	if err != nil {
// 		return "", fmt.Errorf("Error creating GET request: %v", err)
// 	}
// 	// req.Header.Set("User-Agent", fmt.Sprintf("%s-v%s", common.APP_USER_AGENT, common.APP_VERSION))

// 	// reqOut, err := httputil.DumpRequest(req, false)
// 	// fmt.Println(string(reqOut))

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("Error in HTTP request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("Error reading HTTP response body: %v", err)
// 	}

// 	if resp.StatusCode/100 != 2 {
// 		return string(body), fmt.Errorf("HTTP Error %d", resp.StatusCode)
// 	}

// 	return string(body), nil
// }

func purgeCache() {
	concurrentPurgeRoutines <- struct{}{}

	now := time.Now()

	for key, entry := range cache {
		if entry.expiresOn.Before(now) {
			fmt.Println("Purging:", key)
			delete(cache, key)
		}
	}

	<-concurrentPurgeRoutines
}

func getKey(term string, country string, media string) string {
	hash := md5.New()
	io.WriteString(hash, term)
	io.WriteString(hash, country)
	io.WriteString(hash, media)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func main() {
	fmt.Println("*** Welcome to App Store Web Search ***")

	cache = make(map[string]*cacheEntry)
	concurrentPurgeRoutines = make(chan struct{}, 1)

	js.Global().Set("get_app_version", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Invoke(common.APP_VERSION)
		return nil
	}))

	// js.Global().Set("search", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	go func() {
	// 		term := args[0].String()
	// 		country := args[1].String()
	// 		media := args[2].String()
	// 		cb := args[3]

	// 		json, err := search(term, country, "", media, "",
	// 			&http.Client{Timeout: 4 * 1000000000})
	// 		if err != nil {
	// 			cb.Invoke(err.Error(), json)
	// 			return
	// 		} else {
	// 			cb.Invoke(js.Null(), json)
	// 		}
	// 	}()
	// 	return nil
	// }))

	js.Global().Set("get_cache", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		term := args[0].String()
		country := args[1].String()
		media := args[2].String()
		cb := args[3]

		go func() {
			if len(cache) > CACHE_PURGE_THRESHOLD {
				purgeCache()
			}
		}()

		key := getKey(term, country, media)
		if entry := cache[key]; entry != nil {
			cb.Invoke(js.Null(), map[string]interface{}{"key": key, "data": entry.data})
		} else {
			cb.Invoke(fmt.Errorf("cache miss for: %s", key).Error(), nil)
		}

		return nil
	}))

	js.Global().Set("put_cache", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		term := args[0].String()
		country := args[1].String()
		media := args[2].String()
		data := args[3].String()
		cb := args[4]

		key := getKey(term, country, media)
		cache[key] = &cacheEntry{data, time.Now().Add(time.Second * time.Duration(CACHE_TTL))}
		cb.Invoke(js.Null(), key)

		return nil
	}))

	select {}
}
