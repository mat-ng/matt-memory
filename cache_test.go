package mattmemory

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheSetAndGet(t *testing.T) {
	idealItemsPerShard := 10
	loadBalancingInterval := time.Hour
	cache, _ := New(idealItemsPerShard, loadBalancingInterval)

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
	idealItemsPerShard := 10
	loadBalancingInterval := time.Hour
	cache, _ := New(idealItemsPerShard, loadBalancingInterval)

	key := []byte("key1")
	value := []byte("value1")

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
	idealItemsPerShard := 10
	loadBalancingInterval := time.Hour
	cache, _ := New(idealItemsPerShard, loadBalancingInterval)

	key := []byte("key1")
	value := []byte("value1")

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

func TestCacheSetAndGetConcurrently(t *testing.T) {
	idealItemsPerShard := 10
	loadBalancingInterval := time.Hour
	cache, _ := New(idealItemsPerShard, loadBalancingInterval)

	numOps := 10

	var wg sync.WaitGroup
	wg.Add(numOps)

	// Concurrent cache set commands
	for i := 0; i < numOps; i++ {
		go func(idx int) {
			defer wg.Done()
			err := cache.Set([]byte(fmt.Sprintf("key%d", idx)), []byte(fmt.Sprintf("value%d", idx)), 0)
			if err != nil {
				t.Errorf("failed to set key: %v", err)
			}
		}(i)
	}

	wg.Wait()
	wg.Add(numOps)

	// Concurrent cache get commands
	for i := 0; i < numOps; i++ {
		go func(idx int) {
			defer wg.Done()
			retrievedValue, err := cache.Get([]byte(fmt.Sprintf("key%d", idx)))
			if err != nil {
				t.Errorf("failed to get key: %v", err)
			}

			if !bytes.Equal(retrievedValue, []byte(fmt.Sprintf("value%d", idx))) {
				t.Errorf("expected value %s, but got %s", []byte(fmt.Sprintf("value%d", idx)), retrievedValue)
			}
		}(i)
	}

	wg.Wait()
}

func TestCacheLoadBalancing(t *testing.T) {
	idealItemsPerShard := 1
	loadBalancingInterval := 5 * time.Second
	cache, _ := New(idealItemsPerShard, loadBalancingInterval)

	for i := 0; i < 10; i++ {
		err := cache.Set([]byte(fmt.Sprintf("key%d", i)), []byte(fmt.Sprintf("value%d", i)), 10*time.Second)
		if err != nil {
			t.Fatalf("failed to set key: %v", err)
		}
	}

	time.Sleep(7 * time.Second)

	key := []byte("key1")
	value := []byte("value1")

	retrievedValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf("failed to get key: %v", err)
	}

	if !bytes.Equal(retrievedValue, value) {
		t.Fatalf("expected value %s, but got %s", value, retrievedValue)
	}

	time.Sleep(7 * time.Second)

	_, err = cache.Get(key)
	if err == nil {
		t.Fatalf("should have failed to get expired %s", key)
	}
}
