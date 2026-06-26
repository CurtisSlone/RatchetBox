package main

import (
	"golang.org/x/sync/singleflight"
	"time"
)

type Proxy struct {
	cache   *Cache
	fetcher Fetcher
	group   *singleflight.Group
}

func NewProxy(cache *Cache, fetcher Fetcher) *Proxy {
	return &Proxy{
		cache:   cache,
		fetcher: fetcher,
		group:   &singleflight.Group{},
	}
}

func (p *Proxy) Get(key string) (string, error) {
	if value, ok := p.cache.Get(key); ok {
		return value, nil
	}

	result, err, shared := p.group.Do(key, func() (any, error) {
		value, err := p.fetcher.Fetch(key)
		if err != nil {
			return "", err
		}
		p.cache.Put(key, value, 5*time.Minute)
		return value, nil
	})

	if err != nil {
		return "", err
	}

	if shared {
		if value, ok := p.cache.Get(key); ok {
			return value, nil
		}
	}

	return result.(string), nil
}
