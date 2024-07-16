package cache

import (
	"time"
)

type Cacher interface {
	// Retrieves the value associated with a given key from the distributed cache.
	// Returns the value as a byte slice. If they key does not exist, returns an error.
	Get([]byte) ([]byte, error)

	// Stores a value in the cache with the specified key and optional time-to-live (TTL).
	// If the TTL is zero, the key-value pair is held indefinitely.
	// Returns an error if the operation fails.
	Set([]byte, []byte, time.Duration) error

	// Removes the key-value pair of the given key from the distributed cache.
	// If the key does not exist, returns an error.
	Delete([]byte) error

	// Checks if the distributed cache contains a value associated with the given key.
	// If the key-value pair exists in the cache, returns true. Otherwise, returns false.
	Has([]byte) bool
}
