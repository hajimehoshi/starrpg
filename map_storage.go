package starrpg

import (
	"json"
	"os"
	"strconv"
)

type mapStorageImpl struct {
	storage Storage
}

func NewMapStorage(storage Storage) MapStorage {
	return &mapStorageImpl{storage:storage}
}

func (s *mapStorageImpl) Get(key string) (map[string]string, os.Error) {
	bytes := s.storage.Get(key)
	if bytes == nil {
		return nil, nil
	}
	obj := map[string]string{}
	if err := json.Unmarshal(bytes, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *mapStorageImpl) GetWithPrefix(prefix string) (map[string]map[string]string, os.Error) {
	entries := s.storage.GetWithPrefix(prefix) // returns entries with full keys?
	objs := map[string]map[string]string{}
	for key, bytes := range entries {
		obj := map[string]string{}
		if err := json.Unmarshal(bytes, &obj); err != nil {
			return nil, err
		}
		objs[key] = obj
	}
	return objs, nil
}

func (s *mapStorageImpl) Set(key string, obj map[string]string) os.Error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	s.storage.Set(key, bytes)
	return nil
}

func (s *mapStorageImpl) Delete(key string) bool {
	return s.storage.Delete(key)
}

func (s *mapStorageImpl) Update(key string, f func(obj map[string]string) os.Error) os.Error {
	err := s.storage.Update(key, func (bytes []byte) ([]byte, os.Error) {
		obj := map[string]string{}
		if 0 < len(bytes) {
			if err := json.Unmarshal(bytes, &obj); err != nil {
				return nil, err
			}
		}
		if err := f(obj); err != nil {
			return nil, err
		}
		bytes, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	})
	return err
}

func (s *mapStorageImpl) Inc(key, subKey string) (uint64, os.Error) {
	num := uint64(0)
	err := s.Update(key, func (obj map[string]string) os.Error {
		numStr, ok := obj[subKey]
		if ok {
			num2, err := strconv.Atoui64(numStr)
			if err != nil {
				return err
			}
			num = num2
		}
		obj[subKey] = strconv.Uitoa64(num + 1)
		return nil
	})
	if err != nil {
		return 0, err
	}
	return num + 1, nil
}
