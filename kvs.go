package main

import (
	"os"
	"encoding/gob"
	"time"
	"flag"
	"fmt"
	"math/rand"
	"sync"
//	"net/http"
)

var counter = struct{
    sync.RWMutex
    m map[string]int
}{m: make(map[string]int)}

func saveMap(file string, pats *map[string]string) {
	f, err := os.Create(file)
	if err != nil {
			panic("cant open file")
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(pats); err != nil {
			panic("cant encode")
		}
}

func loadMap(file string) (pats *map[string]string) {
	f, err := os.Open(file)
	if err != nil {
		panic("cant open file")
	}
	defer f.Sync()
	defer f.Close()

	enc := gob.NewDecoder(f)
	if err := enc.Decode(&pats); err != nil {
		panic("cant decode")
	}

	return pats
}

func writeValue(m map[string]string, key string, value string) {
	m["key"]="value"
}

func readValue(m map[string]string, key string) (value string) {
	return m["key"]
}

func worker (m *map[string]string) {

}


func main() {
	// create map
	m := make(map[string]string)

	var nKeys int
	var nRndWrites int
	var nRndReads int
	var seed int
	line := "==================================================================="
	flag.IntVar(&nKeys, "keys", 100000, "Number of key / value pairs")
	flag.IntVar(&nRndWrites, "writes", 100000, "Number of random writes")
	flag.IntVar(&nRndReads, "reads", 100000, "Number of random reads")
	flag.IntVar(&seed, "seed", 1337, "seed for rand")
	flag.Parse()

	rand.Seed(int64(seed))

	// start server
//	mux := http.NewServeMux()
/*	mux.Handle("/api/", apiHandler{})
/*	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
*/
	// create n key / values
	start := time.Now()
	fmt.Println(line)
	fmt.Printf("Create %d key / value pairs\n", nKeys)
	for i:=0; i < nKeys; i++ {
		//m["key"+string(i)]="value"+string(i)
		writeValue(m, "key"+string(i), "value"+string(i))
	}
	fmt.Println("Time: ", time.Since(start),
		"time per key: ", time.Since(start).Nanoseconds()/int64(nKeys), "ns")

	// write random value
	start = time.Now()
	fmt.Println(line)
	fmt.Printf("Do %d random writes\n", nRndWrites)
	for i:=0; i < nRndWrites; i++ {
		//m["key"+string(rand.Intn(nKeys))]="random write"+string(i)
		writeValue(m, "key"+string(rand.Intn(nKeys)), "random write"+string(i))
	}
	fmt.Println("Time: ", time.Since(start),
	"time per write: ", time.Since(start).Nanoseconds()/int64(nRndWrites), "ns")

	// read random value
	start = time.Now()
	fmt.Println(line)
	fmt.Printf("Do %d random reads\n", nRndReads)
	tmp := ""
	for i:=0; i < nRndReads; i++ {
		//tmp = m["key"+string(rand.Intn(nKeys))]
		tmp = readValue(m,"key"+string(rand.Intn(nKeys)))
	}
	fmt.Println("Last value: ", tmp)
	fmt.Println("Time: ", time.Since(start),
	"time per read: ", time.Since(start).Nanoseconds()/int64(nRndReads), "ns")

	// write to file
	start = time.Now()
	fmt.Println(line)
	fmt.Println("Write data to file")
	saveMap("map.bin", &m)
	fmt.Println("Time: ", time.Since(start),
	"time per dump: ", time.Since(start).Nanoseconds()/int64(nKeys), "ns")
	fmt.Println(line)
}
