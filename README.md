# Matt Memory ⚡

`matt-memory` is a customizable and scalable distributed cache package for Go, providing rapid data access that can be tailored to meet specific application requirements.

## Flexibility
Users have granular control over key parameters of the cache, so that they can customize the cache's behaviour according to specific needs: 

1.  `idealItemsPerShard`: By specifying the ideal number of items per shard, users can customize the cache to achieve their optimal balance between fast data access and efficient memory usage.

2.  `loadBalancingInterval`: Users can customize the time interval that the cache performs automatic load balancing, giving them the option to optimize workload distribution based on their specific needs.

## Features

🧩 <b>Consistent Hashing:</b> Distributes data across shards using FNV-1a hashing, supporting a balanced distribution across horizontal partitions, while enabling fast cache data retrieval and storage.

📊 <b>Automatic Load Balancing:</b> Redistributes load automatically across shards, optimizing resource usage and ensuring consistent data access speeds.

⌛ <b>TTL Support:</b> Enables the storage of key-value pairs with Time-to-Live (TTL) expirations, allowing for the efficient management of cache memory and ensuring that outdated data is automatically evicted.

🔐 <b>Concurrent Safe:</b> Ensures safe operations in concurrent access scenarios with mutex locks, guaranteeing thread safety by allowing multiple concurrent reads while enforcing exclusive access during writes.

## Usage
`matt-memory` can be installed in your Go project as follows:
```bash
go get github.com/mat-ng/matt-memory
```
Provided below is an example of how to use the cache:
```Go
package main

import (
	"fmt"
	"time"

	mattmemory "github.com/mat-ng/matt-memory"
)

func main() {
	// Define the cache parameters
	idealItemsPerShard := 10
	loadBalancingInterval := time.Hour

	// Create a new cache instance
	cache, err := mattmemory.New(idealItemsPerShard, loadBalancingInterval)
	if err != nil {
		fmt.Printf("Error creating cache: %v\n", err)
		return
	}

	// Set a key-value pair in the cache with a TTL of 30 seconds
	key := []byte("key1")
	value := []byte("value1")
	err = cache.Set(key, value, 30*time.Second)
	if err != nil {
		fmt.Println("error setting value:", err)
		return
	}

	// Get the value from the cache
	result, err := cache.Get(key)
	if err != nil {
		fmt.Println("error getting value:", err)
		return
	} else {
		fmt.Println("value retrieved:", string(result))
	}

	// Check if the key exists in the cache
	if cache.Has(key) {
		fmt.Println("key exists in the cache")
	} else {
		fmt.Println("key does not exist in the cache")
	}

	// Delete the key from the cache
	err = cache.Delete(key)
	if err != nil {
		fmt.Println("error deleting key:", err)
		return
	} else {
		fmt.Println("key deleted successfully")
	}
}
```