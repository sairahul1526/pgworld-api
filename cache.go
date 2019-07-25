package main

import (
	"fmt"

	"github.com/coocood/freecache"
)

// freecache
var cache *freecache.Cache

func getValueFromCache(key string) string {
	got, err := cache.Get([]byte(key))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(got)
}

func setValueInCache(key string, value string) {
	cache.Set([]byte(key), []byte(value), 0)
}

func initcache() {
	cacheSize := 100 * 1024 * 1024 // 100 mb
	cache = freecache.NewCache(cacheSize)
}
