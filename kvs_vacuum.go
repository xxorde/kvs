package kvs

import (
	"time"
)

// Vacuum deletes all dead tuple.
func (s *Kvs) Vacuum() {
	s.Lock()
	defer s.Unlock()
	var dirty int64
	now := time.Now().Unix()

	// Iterate over the unsorted keys.
	for k := range s.values {
		// If key is still valid, leave it.
		if s.values[k].TTL == 0 {
			continue
		} else if s.values[k].TTL > now {
			// Count every non permanent key remaining.
			dirty++
			continue
		} else {
			// If key is end of life, delete it.
			delete(s.values, k)
		}
	}
	s.dirtyUpperLimit = dirty
	s.LastVacuum = time.Now()
	s.VacuumCount++
}

// autoVacuum periodically calls Vacuum() if it estimates that the kvs
// contains too many dead tuple.
func (s *Kvs) autoVacuum() {
	time.Sleep(s.AutoVacuumNaptime)
	if s.dirtyUpperLimit > s.AutoVacuumThreshold &&
		float64(s.dirtyUpperLimit) > (float64(s.Len())*s.AutoVacuumScaleFactor) {
		s.Vacuum()
	}
	go s.autoVacuum()
}
