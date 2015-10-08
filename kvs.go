// 2015 Alexander Sosna <alexander@xxor.de>

// Package kvs is a simple key value store
package kvs

import (
	"sync"
	"time"
)

// Tupel is the type stored in the kvs
type Tupel struct {
	Value string
	TTL   int64
}

// Kvs is the key value store
type Kvs struct {
	sync.RWMutex
	values map[string]Tupel
}

// NewKvs is the constuctor, creating a new Kvs instance
func NewKvs() *Kvs {
	kvs := new(Kvs)
	kvs.values = make(map[string]Tupel)
	return kvs
}

// Len returns the number of stored tuples
func (s *Kvs) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.values)
}

// PutTTL stores a key/value pair with an time to live, ttl.
func (s *Kvs) PutTTL(key string, value string, ttl time.Time) {
	s.Lock()
	defer s.Unlock()
	tmpTupel := Tupel{value, ttl.Unix()}
	s.values[key] = tmpTupel
}

// Put stores a key/value pair.
func (s *Kvs) Put(key string, value string) {
	s.Lock()
	defer s.Unlock()
	tmpTupel := Tupel{value, 0}
	s.values[key] = tmpTupel
}

// Get returns the value for a key
func (s *Kvs) Get(key string) (value string) {
	// RUnlock must be called before the delete, defer not possible
	s.RLock()
	tmpTupel := s.values[key]
	if tmpTupel.TTL == 0 {
		s.RUnlock()
		return tmpTupel.Value
	}

	if tmpTupel.TTL > time.Now().Unix() {
		// Value is still valid
		s.RUnlock()
		return tmpTupel.Value
	}

	// unlock and delete the key
	s.RUnlock()
	s.Delete(key)

	// return empty string, because the value is no longer valid
	return ""
}

// Delete removes value with given key
func (s *Kvs) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.values, key)
}

// Exists tests if given key hast a value
func (s *Kvs) Exists(key string) (exist bool) {
	s.RLock()
	defer s.RUnlock()
	_, exist = s.values[key]
	return exist
}
