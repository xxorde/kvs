// 2015 Alexander Sosna <alexander@xxor.de>

// Package kvs is a simple key value store for concurrent usage.
// It is designed to store values of the type Tuple.
//
// Get a new key value store:
//    store := kvs.NewKvs()
//
// Store a value for key:
//    store.Put("key1", "value1")
//
// Get value for a key:
//    store.Get("key1")
//
// Test if value exists for given key:
//    store.Exists("key")
//
// Output data as yaml:
//    store.Yaml()
//
// Output data as json:
//    store.JSON()
//
// Write data as yaml to an io.Writer:
//    store.DumpYaml(io.Writer)
//
// Write import yaml from an io.Reader:
//    store.ImportYaml(io.Reader)
//
// The functions Yaml() and JSON() will create the export data in memory.
// If you are planing to export much data or have limited memory you should use DumpYaml(io.Writer).
package kvs

import (
	"sync"
	"time"
)

// Tuple is the type stored in the kvs.
// The payload is the string Value.
// With TTL can the lifespan be limited.
// It holds a Unix epoch (seconds), if these time is reached, the payload is invalid.
// If TTL is not set (value 0) the payload is always valid.
type Tuple struct {
	Value string
	TTL   int64
}

// Kvs is the key value store
// It holds an RWMutex to protect the data for concurrent usage and a map.
// The map called values, contains the actual payload as values of type Tuples.
// dirty is an upper boundary on how many non permanent tuples are stored.
type Kvs struct {
	sync.RWMutex
	values                map[string]Tuple
	LastVacuum            time.Time
	VacuumCount           int64
	dirtyUpperLimit       int64
	AutoVacuumEnabled     bool
	AutoVacuumNaptime     time.Duration
	AutoVacuumThreshold   int64
	AutoVacuumScaleFactor float64
}

// NewKvs is the constuctor, creating a new Kvs instance.
func NewKvs() *Kvs {
	kvs := new(Kvs)
	kvs.values = make(map[string]Tuple)
	// Enable autoVacuum
	kvs.AutoVacuumEnabled = true
	// Activate autoVacuum every 10 Minutes
	kvs.AutoVacuumNaptime = 10 * time.Minute
	// Threshold for max tupel with TTL before triggering Vacuum
	kvs.AutoVacuumThreshold = 1000
	// Rate of kvs that could be filled with non permanent tuples, before
	// autoVacuum triggers a Vacuum run.
	kvs.AutoVacuumScaleFactor = 0.2
	// Start autoVacuum (worker process)
	go kvs.autoVacuum()
	return kvs
}

// Len returns the number of stored tuples.
// It is not guarantied that all values are still valid and not end of life.
func (s *Kvs) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.values)
}

// PutTTL stores a key/value pair with an time to live, ttl.
func (s *Kvs) PutTTL(key string, value string, ttl time.Time) {
	s.Lock()
	defer s.Unlock()
	// Create an new tuple, set values and add it to the kvs.
	tmpTuple := Tuple{value, ttl.Unix()}
	s.values[key] = tmpTuple
	// Increment the dirty counter
	s.dirtyUpperLimit++
}

// Put stores a key/value pair.
func (s *Kvs) Put(key string, value string) {
	s.Lock()
	defer s.Unlock()
	tmpTuple := Tuple{value, 0}
	s.values[key] = tmpTuple
}

// Get returns the value for a given key
func (s *Kvs) Get(key string) (value string) {
	// Get a read lock
	s.RLock()
	defer s.RUnlock()
	tmpTuple := s.values[key]

	// If the tuple is valid unlock the mutex and return the value
	if tmpTuple.valid() {
		return tmpTuple.Value
	}

	// The Tuple is not valid, return empty string and delete Tuple (new go routine)
	go s.Delete(key)
	return ""
}

// Delete removes value with given key
func (s *Kvs) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.values, key)
}

// Exists tests if given key hast a value and if it is valid
func (s *Kvs) Exists(key string) (exist bool) {
	s.RLock()
	defer s.RUnlock()
	tmpTupel, exist := s.values[key]
	// The tuple does not exists => return false
	if exist == false || tmpTupel.Value == "" {
		return false
	}
	return tmpTupel.valid()
}

// valid tests if given key hast a valid tuple
func (t *Tuple) valid() (valid bool) {
	// If the tuple has no TTL or it is in the future return true
	if t.TTL == 0 || t.TTL > time.Now().Unix() {
		return true
	}
	// The TTL is in the past, return false
	return false
}
