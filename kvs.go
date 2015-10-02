package kvs

import (
	"sync"
	"time"
)

type Tupel struct {
	Value string
	Ttl time.Time
}

type Kvs struct {
	sync.RWMutex
	M map[string]Tupel
}

func NewKvs() *Kvs {
	kvs := new(Kvs)
	kvs.M = make(map[string]Tupel)
	return kvs
}

/*
func (s *Kvs) init(){
	s.M = make(map[string]string)
}
*/

func (s *Kvs) Len() int {
	return len(s.M)
}

func (s *Kvs) PutTtl(key string, value string, ttl time.Time) {
	s.Lock()
	defer s.Unlock()
	tmpTupel := s.M[key]
	tmpTupel.Value = value
	tmpTupel.Ttl = ttl
	s.M[key] = tmpTupel
}

func (s *Kvs) Put(key string, value string) {
	var ttl time.Time
	s.PutTtl(key, value, ttl)
}

func (s *Kvs) Get(key string) (string) {
	s.RLock()
	defer s.RUnlock()
	tmpTupel := s.M[key]
	tmpTtl := tmpTupel.Ttl
	if tmpTtl.Second() == 0 {
		return tmpTupel.Value
	}
	delta := tmpTtl.Sub(time.Now())
	if(delta.Nanoseconds() > 0) {
		// Value is still valid
		return tmpTupel.Value
	} else {
		// delete the key
		go s.Delete(key)
		// return empty string, because the value is no longer valid
		return ""
	}
}

func (s *Kvs) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.M, key)
}

func (s *Kvs) Exists(key string) (value bool) {
	s.RLock()
	_, value = s.M[key]
	s.RUnlock()
	return value
}
