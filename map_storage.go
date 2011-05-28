package starrpg

import (
	"json"
	"os"
	"strconv"
)

type MapStorage struct {
	storage Storage
}

func NewMapStorage(storage Storage) *MapStorage {
	return &MapStorage{storage:storage}
}

func (s *MapStorage) Get(key string) (map[string]string, os.Error) {
	bytes := s.storage.Get(key)
	if bytes == nil {
		return nil, nil
	}
	obj := map[string]string{}
	if err := json.Unmarshal(bytes, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *MapStorage) GetWithPrefix(prefix string) (map[string]map[string]string, os.Error) {
	entries := s.storage.GetWithPrefix(prefix)
	objs := map[string]map[string]string{}
	for key, bytes := range entries {
		obj := map[string]string{}
		if err := json.Unmarshal(bytes, obj); err != nil {
			return nil, err
		}
		objs[key] = obj
	}
	return objs, nil
}

func (s *MapStorage) Set(key string, obj map[string]string) os.Error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	s.storage.Set(key, bytes)
	return nil
}

func (s *MapStorage) Delete(key string) bool {
	return s.storage.Delete(key)
}

func (s *MapStorage) Inc(key, subKey string) (uint64, os.Error) {
	num := uint64(0)
	err := s.storage.Update(key, func (bytes []byte) ([]byte, os.Error) {
		obj := map[string]string{}
		if bytes != nil {
			if err := json.Unmarshal(bytes, obj); err != nil {
				return nil, err
			}
		}
		numStr, ok := obj[subKey]
		if ok {
			num2, err := strconv.Atoui64(numStr)
			if err != nil {
				return nil, err
			}
			num = num2
		}
		obj[subKey] = strconv.Uitoa64(num + 1)
		bytes, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	})
	if err != nil {
		return 0, err
	}
	return num + 1, nil
}
