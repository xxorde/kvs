package kvs

import (
	"encoding/json"
	"io"
	"bytes"
	"bufio"
	"fmt"
	"strconv"
)

func (s *Kvs) String() (str string) {
	return s.Yaml()
}

func (s *Kvs) Yaml() (yaml string) {
	buf := new(bytes.Buffer)
	s.DumpYaml(buf)
	return buf.String()
}

func (s *Kvs) DumpYaml(w io.Writer){
	yaml := "---"
	s.RLock()
	defer s.RUnlock()
	for key := range s.M {
		value := s.M[key].Value
		ttl := s.M[key].Ttl

		yaml += "\n"+key+": ["+value+","
		if ttl != 0 {
			yaml += strconv.FormatInt(ttl, 10)
		}
		yaml += "]"
		w.Write([]byte(yaml))
		yaml = ""
	}
}

func (t *Tupel) parseYaml(line string) (key string){
	pointer := 0
	value := ""
	ttl := ""

	// remove spaces
	for ; pointer < len(line); pointer++ {
		if line[pointer] == ' ' || line[pointer] == '\t' {
			continue
		} else {
			break
		}
	}

	// get key
	for ; pointer < len(line); pointer++ {
		if line[pointer] == ':' {
			break
		}
		if line[pointer] == '\\' {
			pointer++
			continue
		}
		key += string(line[pointer])
	}

	// find ":"
	for ; pointer < len(line); pointer++ {
		if line[pointer] == ':' {
			pointer++
			break
		}
	}

	// find "["
	for ; pointer < len(line); pointer++ {
		if line[pointer] == '[' {
			pointer++
			break
		}
	}

	// remove spaces
	for ; pointer < len(line); pointer++ {
		if line[pointer] == ' ' || line[pointer] == '\t' {
			continue
		} else {
			break
		}
	}

	// get value
	for ; pointer < len(line); pointer++ {
		if line[pointer] == '\\' {
			pointer++
			continue
		}
		if line[pointer] == ',' {
			break
		}
		value += string(line[pointer])
	}
	t.Value = value

	// remove spaces
	for ; pointer < len(line); pointer++ {
		if line[pointer] == ' ' || line[pointer] == ' ' {
			continue
		} else {
			break
		}
	}

	// get ttl
	for ; pointer < len(line); pointer++ {
		if line[pointer] == '\\' {
			pointer++
			continue
		}
		if line[pointer] == ',' {
			break
		}
		ttl += string(line[pointer])
	}

	if ttl != "" {
		ttli, err := strconv.ParseInt(ttl, 10, 32)
		if err != nil {
			panic(err)
		}
		t.Ttl = ttli
	}
	return key
}

func (s *Kvs) ImportYaml(r io.Reader){
//	var noTime time.Time
//	yaml := "---"
	s.Lock()
	defer s.Unlock()
	var newStore Kvs
	newStore = *NewKvs()
	scanner := bufio.NewScanner(r)
	header := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			header = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if header != true {
		panic("No yaml file, no ---")
	}

	for scanner.Scan() {
		line := scanner.Text()
		var tmpTupel Tupel
		key := tmpTupel.parseYaml(line)
		if key != "" {
			s.M[key] = tmpTupel
		}
		fmt.Println("k:",key,"v:",tmpTupel.Value,"ttl:",tmpTupel.Ttl)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	s = &newStore
}

func (s *Kvs) Json() (string){
	s.RLock()
	defer s.RUnlock()
	b, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		panic("error json.MarshalIndent(s)")
	}
	return string(b)
}
