package kvs

import (
	"sync"
	"time"
)

type Tupel struct {
	Value string
	Ttl   int64
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
	tmpTupel.Ttl = ttl.Unix()
	s.M[key] = tmpTupel
}

func (s *Kvs) Put(key string, value string) {
	s.Lock()
	defer s.Unlock()
	tmpTupel := s.M[key]
	tmpTupel.Value = value
	s.M[key] = tmpTupel
}

func (s *Kvs) Get(key string) string {
	s.RLock()
	tmpTupel := s.M[key]
	tmpTtl := tmpTupel.Ttl
	if tmpTtl == 0 {
		s.RUnlock()
		return tmpTupel.Value
	}

	if tmpTtl > time.Now().Unix() {
		// Value is still valid
		s.RUnlock()
		return tmpTupel.Value
	}

	s.RUnlock()

	// delete the key
	go s.Delete(key)

	// return empty string, because the value is no longer valid
	return ""
}

func (s *Kvs) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.M, key)
}

func (s *Kvs) Exists(key string) (exist bool) {
	s.RLock()
	_, exist = s.M[key]
	s.RUnlock()
	return exist
}
