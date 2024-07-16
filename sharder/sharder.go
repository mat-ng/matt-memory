package shard

import "time"

type Sharder interface {
	// Retrieves the value associated with the specified key stored in the shard.
	// Returns the value as a byte slice. If they key does not exist, returns an error.
	Get([]byte) ([]byte, error)

	// Stores a value in the shard with the specified key and time-to-live (TTL).
	// If the TTL is zero, the key-value pair is held indefinitely.
	// Returns an error if the operation fails.
	Set([]byte, []byte, time.Duration) error

	// Removes the key-value pair of the given key stored in the shard.
	// If the key does not exist, returns an error.
	Delete([]byte) error

	// Checks if the shard contains a value associated with the given key.
	// If the key-value pair exists in the shard, returns true. Otherwise, returns false.
	Has([]byte) bool

	// Returns a list of keys stored in the shard.
	Range() []string
}
