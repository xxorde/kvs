package kvs

import (
	"compress/gzip"
	"encoding/gob"
	"io"
	"os"
	"sync"
)

type Kvs struct {
	sync.RWMutex
	M map[string]string
}

func NewKvs() *Kvs {
	kvs := new(Kvs)
	kvs.M = make(map[string]string)
	return kvs
}

/*
func (s *Kvs) init(){
	s.M = make(map[string]string)
}
*/

func (s *Kvs) Len() int {
	return len(s.M)
}

func (s *Kvs) Put(key string, value string) {
	s.Lock()
	s.M[key] = value
	s.Unlock()
}

func (s *Kvs) Get(key string) (value string) {
	s.RLock()
	value = s.M[key]
	s.RUnlock()
	return value
}

func (s *Kvs) Delete(key string) {
	s.Lock()
	delete(s.M, key)
	s.Unlock()
}

func (s *Kvs) Exists(key string) (value bool) {
	s.RLock()
	_, value = s.M[key]
	s.RUnlock()
	return value
}

func (s *Kvs) BackupBin(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic("cant open file")
	}
	defer f.Sync()
	defer f.Close()

	s.ExportBin(f)
}

func (s *Kvs) BackupBinGz(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic("cant open file")
	}
	defer f.Sync()
	defer f.Close()

	z := gzip.NewWriter(f)
	s.ExportBin(z)
}

func (s *Kvs) ExportBin(f io.Writer) {
	s.Lock()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(s.M); err != nil {
		panic("cant encode")
	}
	s.Unlock()
}

func (s *Kvs) RestoreBin(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic("cant open file")
	}
	defer f.Close()
	s.ImportBin(f)
}

func (s *Kvs) RestoreBinGz(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic("cant open file")
	}
	defer f.Close()

	z, err := gzip.NewReader(f)
	if err != nil {
		panic("cant open file")
	}
	s.ImportBin(z)
}

func (s *Kvs) ImportBin(f io.Reader) {
	s.Lock()
	enc := gob.NewDecoder(f)
	if err := enc.Decode(&s.M); err != nil {
		panic("cant decode")
	}
	s.Unlock()
}
