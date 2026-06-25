package main

import "sync"

type Mapper struct {
	mapping map[string]string
	lock    sync.Mutex
}

var urlMapper = newMapper()

func newMapper() Mapper {
	return Mapper{
		mapping: make(map[string]string),
		lock:    sync.Mutex{},
	}
}

func insertMapping(key, originalURL string) {
	urlMapper.lock.Lock()
	defer urlMapper.lock.Unlock()

	urlMapper.mapping[key] = originalURL
}

func fetchMapping(key string) (string, bool) {
	urlMapper.lock.Lock()
	defer urlMapper.lock.Unlock()

	originalURL, exists := urlMapper.mapping[key]
	return originalURL, exists
}
