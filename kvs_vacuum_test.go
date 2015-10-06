package kvs

/*
import (
	"testing"
	"time"
)


func TestVacuum(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		// get new kvs
		store = *NewKvs()

		// set wait time
		waitInSeconds := 1
		waitTime := time.Duration(waitInSeconds) * time.Second

		// set ttl to n seconds
		ttl := time.Now()
		ttl = ttl.Add(waitTime)

		// store values with ttl
		for _, c := range tc.kvs {
			store.PutTTL(c.key, c.value, ttl)
		}

		// the values should still be there!
		storeLen := store.Len()
		if storeLen != len(tc.kvs) {
			t.Errorf("storeLen != len(tc.kvs), %q != %q", storeLen, len(tc.kvs))
		}

		// sleep n secods
		time.Sleep(waitTime)

		// test if dead tuples are still there
		storeLenAfterTTL := store.Len()
		if storeLen != len(tc.kvs) {
			t.Errorf("storeLenAfterTTL != len(tc.kvs), %q != %q", storeLenAfterTTL, len(tc.kvs))
		}

		// clean dead tuples
		store.Vacuum()

		// test if dead tuples are cleaned
		storeLenAfterVacuum := store.Len()
		if storeLenAfterVacuum > 0 {
			t.Errorf("storeLenAfterVacuum > 0, %q", storeLenAfterVacuum)
		}
	}
}
*/
