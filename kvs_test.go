package kvs

import (
	"testing"
	"os"
	"time"
	"fmt"
//	"math/rand"
//	"strconv"
)

var (
	store Kvs
	cases []struct {key, value string}
)

func TestMain(m *testing.M) {
	// create global kvs
	store = *NewKvs()

	// just so we use fmt
	fmt.Println("Test started...")

	// global testcases
	cases = []struct {key, value string}{
		{"key","value"},
		{"Hello", "Welt"},
		{"Hack", "the Planet"},
		{"Hack2", "the Planetthe Planetthe Planetthe Planetthe Planetthe Planet"},
	}
	os.Exit(m.Run())
}

func TestPutGet(t *testing.T) {
	for _, c := range cases {
		store.Put(c.key,c.value)
	}
	tmp := ""
	for _, c := range cases {
		tmp = store.Get(c.key)
		if c.value != tmp {
			t.Errorf("Get(%q) == %q, should be: %q", c.key, tmp, c.value)
		}
	}
}

func TestPutGetTtl(t *testing.T) {
	// set wait time
	waitInSeconds := 1
	waitTime := time.Duration(waitInSeconds)*time.Second

	// set ttl to n seconds
	ttl := time.Now()
	ttl = ttl.Add(waitTime)

	// store values with ttl
	for _, c := range cases {
		store.PutTtl(c.key,c.value,ttl)
	}

	// the values should still be there!
	tmp := ""
	for _, c := range cases {
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
	for _, c := range cases {
		tmp = store.Get(c.key)
		if tmp != empty  {
			t.Errorf("Get(%q) == %q, should be empty", c.key, tmp)
		}
	}
}
