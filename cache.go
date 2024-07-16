package mattmemory

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	shard "github.com/mat-ng/matt-memory/sharder"
)

type Cache struct {
	shards                []*shard.Shard
	idealItemsPerShard    int
	loadBalancingInterval time.Duration

	mutex sync.Mutex
}

func New(idealItemsPerShard int, loadBalancingInterval time.Duration) (*Cache, error) {
	if idealItemsPerShard <= 0 {
		return nil, fmt.Errorf("invalid ideal number of items per shard: %d", idealItemsPerShard)
	}
	if loadBalancingInterval < (5 * time.Second) {
		return nil, fmt.Errorf("invalid load balancing interval (minimum 5 seconds): %d", loadBalancingInterval)
	}

	cache := &Cache{
		shards:                []*shard.Shard{shard.New()},
		idealItemsPerShard:    idealItemsPerShard,
		loadBalancingInterval: loadBalancingInterval,
	}

	// Start load balancing goroutine
	go cache.loadBalance()

	return cache, nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	shard := c.getShardForKey(key)
	return shard.Get(key)
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) error {
	shard := c.getShardForKey(key)
	return shard.Set(key, value, ttl)
}

func (c *Cache) Delete(key []byte) error {
	shard := c.getShardForKey(key)
	return shard.Delete(key)
}

func (c *Cache) Has(key []byte) bool {
	shard := c.getShardForKey(key)
	return shard.Has(key)
}

func (c *Cache) getShardForKey(key []byte) *shard.Shard {
	keyStr := string(key)
	hash := c.hashFunction(keyStr)
	shardIndex := hash % uint32(len(c.shards))

	return c.shards[shardIndex]
}

func (c *Cache) hashFunction(key string) uint32 {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))

	return hasher.Sum32()
}

func (c *Cache) loadBalance() {
	ticker := time.NewTicker(c.loadBalancingInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		c.mutex.Lock()

		// Cache must have at least 1 shard
		idealShardCount := 1
		currentShardCount := len(c.shards)

		// Calculate total item count across all shards
		totalItems := 0
		for _, shard := range c.shards {
			totalItems += len(shard.Range())
		}

		// Calculate desired shard count rounded up to the nearest integer
		if totalItems > 0 {
			idealShardCount = (totalItems + c.idealItemsPerShard - 1) / c.idealItemsPerShard
		}

		// Adjust shard count
		if idealShardCount > currentShardCount {
			c.addShards(idealShardCount - currentShardCount)
		} else if idealShardCount < currentShardCount {
			c.removeShards(currentShardCount - idealShardCount)
		}

		c.mutex.Unlock()
	}
}

func (c *Cache) addShards(n int) error {
	if n <= 0 {
		return fmt.Errorf("invalid number of shards to add: %d", n)
	}

	// Extract items from all shards
	dataToRedistribute, ttlsToRedistribute := extractItemsFromShards(c.shards)

	// Create new slice of shards
	shardsNewLength := len(c.shards) + n
	shardsNew := make([]*shard.Shard, shardsNewLength)
	for i := 0; i < shardsNewLength; i++ {
		shardsNew[i] = shard.New()
	}

	c.shards = shardsNew

	// Redistribute items to new shards
	c.distributeItemsToShards(dataToRedistribute, ttlsToRedistribute)

	return nil
}

func (c *Cache) removeShards(n int) error {
	if n <= 0 || n >= len(c.shards) {
		return fmt.Errorf("invalid number of shards to remove: %d", n)
	}

	// Extract items from shards that will be removed
	dataToRedistribute, ttlsToRedistribute := extractItemsFromShards(c.shards[:n])

	// Remove shards
	c.shards = c.shards[n:]

	// Redistribute items from removed shards
	c.distributeItemsToShards(dataToRedistribute, ttlsToRedistribute)

	return nil
}

func (c *Cache) distributeItemsToShards(dataToDistribute map[string][]byte, ttlsToDistribute map[string]time.Time) {
	for key, value := range dataToDistribute {
		ttl := time.Duration(0)
		if expiration, ok := ttlsToDistribute[key]; ok {
			ttl = time.Until(expiration)
		}
		if ttl < 0 {
			continue // Skip expired items
		}
		c.Set([]byte(key), value, ttl)
	}
}

func extractItemsFromShards(shards []*shard.Shard) (map[string][]byte, map[string]time.Time) {
	dataExtracted := make(map[string][]byte)
	ttlsExtracted := make(map[string]time.Time)

	for _, shard := range shards {
		for _, key := range shard.Range() {
			value, _ := shard.Get([]byte(key))
			dataExtracted[key] = value

			ttl, err := shard.GetTtl([]byte(key))
			if err == nil {
				ttlsExtracted[key] = ttl
			}
		}
	}

	return dataExtracted, ttlsExtracted
}
