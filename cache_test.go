package cache

import (
	"bytes"
	"testing"
	"time"
)

func TestCacheSetAndGet(t *testing.T) {
	cache := New(10)

	key := []byte("key1")
	value := []byte("value1")

	err := cache.Set(key, value, time.Minute)
	if err != nil {
		t.Fatalf("failed to set key: %v", err)
	}

	retrievedValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf("failed to get key: %v", err)
	}

	if !bytes.Equal(retrievedValue, value) {
		t.Fatalf("expected value %s, but got %s", value, retrievedValue)
	}
}

func TestCacheDelete(t *testing.T) {
	cache := New(10)

	key := []byte("key2")
	value := []byte("value2")

	err := cache.Set(key, value, time.Minute)
	if err != nil {
		t.Fatalf("failed to set key: %v", err)
	}

	err = cache.Delete(key)
	if err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}

	_, err = cache.Get(key)
	if err == nil {
		t.Fatalf("expected error when getting deleted key, but got nil")
	}
}

func TestCacheHas(t *testing.T) {
	cache := New(10)

	key := []byte("key3")
	value := []byte("value3")

	err := cache.Set(key, value, time.Minute)
	if err != nil {
		t.Fatalf("failed to set key: %v", err)
	}

	if !cache.Has(key) {
		t.Fatalf("expected cache to have key %s but it does not", key)
	}

	err = cache.Delete(key)
	if err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}

	if cache.Has(key) {
		t.Fatalf("expected cache to not have key %s but it does", key)
	}
}
