package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	const inter = 3 * time.Second
	cache := NewCache(inter)

	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "https://example1.com",
			value: []byte("some sample data"),
		},
		{
			key:   "http://example2.com",
			value: []byte("some more sample data"),
		},
	}

	for _, c := range cases {
		cache.Add(c.key, c.value)
	}

	t.Run("testing addition into cache", func(t *testing.T) {
		for k := range cache.entries {
			if k == cases[0].key {
				return
			}
		}
		t.Errorf("expected to find value %s", cases[0].key)
		return
	})
}

func TestGet(t *testing.T) {
	const inter = 3 * time.Second
	cache := NewCache(inter)

	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "https://example1.com",
			value: []byte("some sample data"),
		},
		{
			key:   "http://example2.com",
			value: []byte("some more sample data"),
		},
	}

	for _, c := range cases {
		cache.Add(c.key, c.value)
	}

	t.Run(fmt.Sprintf("testing to find %s", cases[1].value), func(t *testing.T) {
		if _, ok := cache.Get(cases[1].key); !ok {
			t.Errorf("expected to find value %s", cases[0].key)
			return
		}
	})
}

func TestReadLoop(t *testing.T) {
	const inter = 5 * time.Millisecond
	const waitTime = inter + 10*time.Millisecond
	const testKey = "http://example.com"
	cache := NewCache(inter)

	cache.Add(testKey, []byte("some amount of data"))
	if _, ok := cache.Get(testKey); !ok {
		t.Errorf("Couldn't find key: %s", testKey)
		return
	}

	time.Sleep(waitTime)

	if _, ok := cache.Get(testKey); ok {
		t.Errorf("Found key: %s after it should have been deleted.", testKey)
		return
	}
}
