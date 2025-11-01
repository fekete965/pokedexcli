package internal

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = time.Second * 5

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://test-key.com",
			val: []byte("test-data"),
		},
		{
			key: "https://test-key.com/path",
			val: []byte("additional-test-data"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Text case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)

			val, ok := cache.Get(c.key)

			if !ok {
				t.Errorf("expected to find key: %v", c.key)
				return
			}

			if string(val) != string(c.val) {
				t.Errorf("expected to find the correct value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = time.Millisecond * 5
	const waitTime = baseTime + time.Millisecond*5

	testCacheKey := "https://test-key.com"

	cache := NewCache(baseTime)
	cache.Add(testCacheKey, []byte("testdata"))

	_, ok := cache.Get(testCacheKey)
	if !ok {
		t.Errorf("expected to find key: %v", testCacheKey)
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(testCacheKey)
	if ok {
		t.Errorf("expected to not find key: %v", testCacheKey)
		return
	}
}
