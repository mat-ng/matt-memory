package shard

import (
	"fmt"
	"sync"
	"time"
)

type Shard struct {
	data  map[string][]byte
	ttl   map[string]time.Time
	Mutex sync.RWMutex
}

func New() *Shard {
	return &Shard{
		data: make(map[string][]byte),
		ttl:  make(map[string]time.Time),
	}
}

func (s *Shard) Get(key []byte) ([]byte, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	keyStr := string(key)

	val, ok := s.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyStr)
	}

	return val, nil
}

func (s *Shard) Set(key []byte, value []byte, ttl time.Duration) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	keyStr := string(key)

	s.data[keyStr] = value

	if ttl > 0 {
		s.ttl[keyStr] = time.Now().Add(ttl)
		go func() {
			<-time.After(ttl)
			delete(s.data, keyStr)
			delete(s.ttl, keyStr)
		}()
	}

	return nil
}

func (s *Shard) Delete(key []byte) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	keyStr := string(key)

	delete(s.data, keyStr)
	delete(s.ttl, keyStr)

	return nil
}

func (s *Shard) Has(key []byte) bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	_, ok := s.data[string(key)]

	return ok
}

func (s *Shard) GetTtl(key []byte) (time.Time, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	keyStr := string(key)

	if !s.Has(key) {
		return time.Time{}, fmt.Errorf("key (%s) not found", keyStr)
	}

	ttl, ok := s.ttl[keyStr]
	if !ok {
		return time.Time{}, fmt.Errorf("key (%s) does not have a ttl", keyStr)
	}

	return ttl, nil
}

func (s *Shard) Range() []string {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	keys := make([]string, 0, len(s.data))

	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}
