package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mu       *sync.Mutex
}

func NewCache(inter time.Duration) Cache {
	ch := Cache{
		entries:  map[string]cacheEntry{},
		interval: inter,
		mu:       &sync.Mutex{},
	}

	go ch.readLoop()

	return ch
}

func (ch *Cache) Add(k string, val []byte) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.entries[k] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (ch *Cache) Get(k string) ([]byte, bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	entry, ok := ch.entries[k]
	if !ok {
		return nil, ok
	}

	return entry.val, ok
}

func (ch *Cache) readLoop() {
	tiker := time.NewTicker(ch.interval)
	defer tiker.Stop()

	for {
		select {
		case t := <-tiker.C:
			ch.mu.Lock()
			for k, entry := range ch.entries {
				diff := t.Sub(entry.createdAt)
				if diff > ch.interval {
					delete(ch.entries, k)
				}
			}
			ch.mu.Unlock()
		}
	}
}
