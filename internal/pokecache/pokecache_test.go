package pokecache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(5 * time.Second)
	if cache == nil {
		t.Errorf("NewCache() returned nil")
	}
}
