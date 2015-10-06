package kvs

import (
	"testing"
	"os"
	"time"
	"fmt"
	"bytes"
//	"math/rand"
//	"strconv"
)

type testtupel struct {
	key string
	value string
	ttl int64
}

type testcase struct {
	kvs []testtupel
	yaml string
	json string
}

var (
	store Kvs
	cases []testcase
	line string
)

func TestMain(m *testing.M) {
	// create global kvs

	// just so we use fmt
	fmt.Println("Test started...")

	// global testtc01Kvs
	cases = []testcase{
		testcase{
			[]testtupel{
				{"key", "value", 0},
				{"Hello", "Welt", 0},
				{"Hack", "the Planet", 0},
				{"Hack2", "Planetthe Planetthe Planet", 0},
			},
			`---
key: [value,]
Hello: [Welt,]
Hack: [the Planet,]
Hack2: [Planetthe Planetthe Planet,]`,
			"{\n \"M\": {\n  \"Hack\": {\n   \"Value\": \"the Planet\",\n   \"Ttl\": 0\n  },\n  \"Hack2\": {\n   \"Value\": \"Planetthe Planetthe Planet\",\n   \"Ttl\": 0\n  },\n  \"Hello\": {\n   \"Value\": \"Welt\",\n   \"Ttl\": 0\n  },\n  \"key\": {\n   \"Value\": \"value\",\n   \"Ttl\": 0\n  }\n }\n}",
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
			"{\n \"M\": {\n  \"1\": {\n   \"Value\": \"111\",\n   \"Ttl\": 0\n  },\n  \"2\": {\n   \"Value\": \"222\",\n   \"Ttl\": 0\n  },\n  \"3\": {\n   \"Value\": \"333\",\n   \"Ttl\": 0\n  },\n  \"4\": {\n   \"Value\": \"444\",\n   \"Ttl\": 0\n  }\n }\n}",
		},
	}

	line = "===================================="
	os.Exit(m.Run())
}

func TestPutGet(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		// get new kvs
		store = *NewKvs()

		for _, c := range tc.kvs {
			store.Put(c.key,c.value)
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

func TestPutGetTtl(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		// get new kvs
		store = *NewKvs()

		// set wait time
		waitInSeconds := 1
		waitTime := time.Duration(waitInSeconds)*time.Second

		// set ttl to n seconds
		ttl := time.Now()
		ttl = ttl.Add(waitTime)

		// store values with ttl
		for _, c := range tc.kvs {
			store.PutTtl(c.key,c.value,ttl)
		}

		//fmt.Println(store.Json())
		//fmt.Println(store.Yaml())

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
			if tmp != empty  {
				t.Errorf("Get(%q) == %q, should be empty", c.key, tmp)
			}
		}
	}
}

func TestJson(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		// get new kvs
		store = *NewKvs()

		// insert testcase
		for _, c := range tc.kvs {
			store.Put(c.key,c.value)
		}

		// test if json matches testcase
		if store.Json() != tc.json {
			t.Errorf("json does not match testcase:\n%q\n !=\n%q", store.Json(), tc.json)
		}

		//fmt.Println(store.Json())
	}
}

func TestYaml(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		store := *NewKvs()
		for _, c := range tc.kvs {
			store.Put(c.key,c.value)
		}
		if store.Yaml() != tc.yaml {
			t.Errorf("yaml does not match testcase:\n%q\n !=\n%q", store.Yaml(), tc.yaml)
		}
	}
}

func TestDumpYaml(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		store := *NewKvs()
		for _, c := range tc.kvs {
			store.Put(c.key,c.value)
		}
		buf := new(bytes.Buffer)
		store.DumpYaml(buf)
		if buf.String() != tc.yaml {
			t.Errorf("DumpYaml does not match testcase:\n%q\n !=\n%q", buf, tc.yaml)
		}
	}
}
/*
func TestImportYaml(t *testing.T) {
	fmt.Println("TestImportYaml")
	store = *NewKvs()
	dump := new(bytes.Buffer)
	TestPutGet(t)
	fmt.Println()
	store.DumpYaml(dump)

	store = *NewKvs()
	store.ImportYaml(dump)
	fmt.Println()
	store.DumpYaml(os.Stdout)
	fmt.Println(line)
}
*/
