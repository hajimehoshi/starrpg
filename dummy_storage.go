package starrpg

import (
	"os"
	"strings"
)

type DummyStorage map[string][]byte

func (s *DummyStorage) Get(key string) []byte {
	bytes, ok := (*s)[key]
	if !ok {
		return nil
	}
	return bytes
}

func (s *DummyStorage) GetWithPrefix(prefix string) map[string][]byte {
	entries := map[string][]byte{}
	for key, value := range *s {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		entries[key] = value
	}
	return entries
}

func (s *DummyStorage) Set(key string, value []byte) {
	(*s)[key] = value
}

func (s *DummyStorage) Delete(key string) bool {
	if _, ok := (*s)[key]; !ok {
		return false
	}
	(*s)[key] = nil, false
	return true
}

func (s *DummyStorage) Update(key string, f func([]byte) ([]byte, os.Error)) os.Error {
	bytes := (*s)[key]
	oldBytes := make([]byte, len(bytes))
	copy(oldBytes, bytes)
	newBytes, err := f(oldBytes)
	if err != nil {
		return err
	}
	if newBytes != nil {
		(*s)[key] = newBytes
	} else {
		(*s)[key] = nil, false
	}
	return nil
}
