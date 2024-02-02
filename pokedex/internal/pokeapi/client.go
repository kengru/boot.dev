package pokeapi

import (
	"net/http"
	"time"

	"internal/pokecache"
)

type PokeClient struct {
	cache      pokecache.Cache
	httpClient http.Client
}

func NewPokeClient(timeout, cacheInterval time.Duration) PokeClient {
	return PokeClient{
		cache: pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
