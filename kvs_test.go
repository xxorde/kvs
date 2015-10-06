package kvs

import (
	"os"
	"sync"
	"testing"
	"time"
)

type testtupel struct {
	key   string
	value string
	ttl   int64
}

type testcase struct {
	kvs  []testtupel
	yaml string
	json string
}

var (
	store Kvs
	cases []testcase
	line  string
)

func TestMain(m *testing.M) {
	// global testcases
	cases = []testcase{
		testcase{
			[]testtupel{
				{"key", "value", 0},
				{"Hello", "Welt", 0},
				{"Hack", "the Planet", 0},
				{"Hack2", "Planetthe Planetthe Planet", 0},
			},
			`---
Hack: [the Planet,]
Hack2: [Planetthe Planetthe Planet,]
Hello: [Welt,]
key: [value,]`,
			"{\n \"M\": {\n  \"Hack\": {\n   \"Value\": \"the Planet\",\n   \"TTL\": 0\n  },\n  \"Hack2\": {\n   \"Value\": \"Planetthe Planetthe Planet\",\n   \"TTL\": 0\n  },\n  \"Hello\": {\n   \"Value\": \"Welt\",\n   \"TTL\": 0\n  },\n  \"key\": {\n   \"Value\": \"value\",\n   \"TTL\": 0\n  }\n }\n}",
		},
		testcase{
			[]testtupel{
				{"1", "111", 0},
				{"2", "222", 0},
				{"3", "333", 0},
				{"4", "444", 0},
			},
			`---
1: [111,]
2: [222,]
3: [333,]
4: [444,]`,
			"{\n \"M\": {\n  \"1\": {\n   \"Value\": \"111\",\n   \"TTL\": 0\n  },\n  \"2\": {\n   \"Value\": \"222\",\n   \"TTL\": 0\n  },\n  \"3\": {\n   \"Value\": \"333\",\n   \"TTL\": 0\n  },\n  \"4\": {\n   \"Value\": \"444\",\n   \"TTL\": 0\n  }\n }\n}",
		},
	}

	os.Exit(m.Run())
}

func TestPutGet(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		// get new kvs
		store = *NewKvs()

		for _, c := range tc.kvs {
			store.Put(c.key, c.value)
		}
		tmp := ""
		for _, c := range tc.kvs {
			tmp = store.Get(c.key)
			if c.value != tmp {
				t.Errorf("Get(%q) == %q, should be: %q", c.key, tmp, c.value)
			}
		}
	}
}

func TestPutGetTTL(t *testing.T) {
	// sync to use go routines for testing
	var wg sync.WaitGroup

	// range over testcases
	for _, tc := range cases {
		var tmpTc testcase
		tmpTc = tc
		wg.Add(1)
		go func(tc testcase) {
			defer wg.Done()
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
			tmp := ""
			for _, c := range tc.kvs {
				tmp = store.Get(c.key)
				if c.value != tmp {
					t.Errorf("Get(%q) == %q, should be: %q", c.key, tmp, c.value)
				}
			}

			// sleep n secods
			time.Sleep(waitTime)

			// test if the values are absent!
			empty := ""
			tmp = ""
			for _, c := range tc.kvs {
				tmp = store.Get(c.key)
				if tmp != empty {
					t.Errorf("Get(%q) == %q, should be empty", c.key, tmp)
				}
			}
		}(tmpTc)
	}

	// wait for all go routines
	wg.Wait()
}
