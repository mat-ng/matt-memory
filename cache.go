package cache

import (
	"hash/fnv"
	"time"

	shard "github.com/mat-ng/matt-memory/sharder"
)

type Cache struct {
	shards []*shard.Shard
}

func New(shardCount int) *Cache {
	shards := make([]*shard.Shard, shardCount)

	for i := 0; i < shardCount; i++ {
		shards[i] = shard.New()
	}

	cache := &Cache{
		shards: shards,
	}

	return cache
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
