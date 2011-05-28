package starrpg

import (
	"strconv"
	"strings"
)

type DummyStorage map[string][]byte

func (s *DummyStorage) Get(key string) []byte {
	item, ok := (*s)[key]
	if !ok {
		return nil
	}
	return item
}

func (s *DummyStorage) GetWithPrefix(prefix string) []*StorageEntry {
	entries := make([]*StorageEntry, 0)
	for key, value := range *s {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		entry := &StorageEntry{Key:key, Value:value}
		entries = append(entries, entry)
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

func (s *DummyStorage) Inc(key string) (uint64, bool) {
	value, ok := (*s)[key]
	if !ok {
		(*s)[key] = []byte("1")
		return 1, true
	}
	numValue, err := strconv.Atoui64(string(value))
	if err != nil {
		return 0, false
	}
	(*s)[key] = []byte(strconv.Uitoa64(numValue + 1))
	return numValue + 1, true
}
