package kvs

import (
	"bytes"
	"testing"
)

func TestJSON(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		// check if json sample is available
		if tc.json == "" {
			continue
		}

		// get new kvs
		store := *NewKvs()

		// insert testcase
		for _, c := range tc.kvs {
			store.Put(c.key, c.value)
		}

		// test if JSON matches testcase
		if store.JSON() != tc.json {
			t.Errorf("json does not match testcase:\n%q\n !=\n%q", store.JSON(), tc.json)
		}

		//fmt.Println(store.JSON())
	}
}

func TestYaml(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		store := *NewKvs()
		for _, c := range tc.kvs {
			store.Put(c.key, c.value)
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
			store.Put(c.key, c.value)
		}
		buf := new(bytes.Buffer)
		store.DumpYaml(buf)
		if buf.String() != tc.yaml {
			t.Errorf("DumpYaml does not match testcase:\n%q\n !=\n%q", buf, tc.yaml)
		}
	}
}

func TestImportYaml(t *testing.T) {
	// range over testcases
	for _, tc := range cases {
		store := *NewKvs()
		dump := new(bytes.Buffer)
		// fill kvs with testcase
		for _, c := range tc.kvs {
			store.Put(c.key, c.value)
		}
		// dump data in dump
		store.DumpYaml(dump)

		// initialize new kvs
		store = *NewKvs()

		// import dump
		store.ImportYaml(dump)

		// test if all tuples are present
		tmp := ""
		for _, c := range tc.kvs {
			tmp = store.Get(c.key)
			if c.value != tmp {
				t.Errorf("Get(%q) == %q, should be: %s", c.key, tmp, c.value)
			}
		}
	}
}
