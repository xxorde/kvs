package kvs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"sort"
	"strconv"
)

func (s *Kvs) String() (str string) {
	return s.Yaml()
}

// Yaml exports kvs as simple yaml file.
func (s *Kvs) Yaml() (yaml string) {
	buf := new(bytes.Buffer)
	s.DumpYaml(buf)
	return buf.String()
}

// DumpYaml writes kvs as simple yaml to a io.Writer
// Expired (invalid) keys are not automatically removed!
func (s *Kvs) DumpYaml(w io.Writer) {
	yaml := "---"
	s.RLock()
	defer s.RUnlock()

	// get all keys and sort them
	var keys []string
	for k := range s.values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// iterate over the sorted keys
	for _, k := range keys {
		value := s.values[k].Value
		ttl := s.values[k].TTL

		yaml += "\n" + k + ": [" + value + ","
		if ttl != 0 {
			yaml += strconv.FormatInt(ttl, 10)
		}
		yaml += "]"
		w.Write([]byte(yaml))
		yaml = ""
	}
}

func (t *Tuple) parseYaml(line string) (key string) {
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
		if line[pointer] == ' ' || line[pointer] == '\t' {
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
		t.TTL = ttli
	}
	return key
}

// ImportYaml reads kvs as simple yaml from a io.Reader
func (s *Kvs) ImportYaml(r io.Reader) {
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
		var tmpTuple Tuple
		key := tmpTuple.parseYaml(line)
		if key != "" {
			s.values[key] = tmpTuple
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	s = &newStore
}

// JSON exports kvs as json file.
func (s *Kvs) JSON() (jsonOutput string) {
	s.RLock()
	defer s.RUnlock()
	bytes, err := json.MarshalIndent(s.values, "", " ")
	if err != nil {
		panic("error json.valuesarshalIndent(s)")
	}
	return string(bytes)
}
