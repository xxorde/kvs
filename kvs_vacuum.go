package kvs

import (
	"time"
)

// Vacuum cleans all dead tuples
func (s *Kvs) Vacuum() {
	s.Lock()
	defer s.Unlock()
	now := time.Now().Unix()

	// iterate over the unsorted keys
	for k := range s.values {
		// if key is still valid, leave it
		if s.values[k].TTL == 0 || s.values[k].TTL > now {
			continue
		} else {
			// if key is end of life, delete it
			delete(s.values, k)
		}
	}
}
