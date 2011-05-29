package starrpg

import (
	"os"
	"strings"
	"strconv"
)

type MapStorage interface {
	Get(key string) (map[string]string, os.Error)
	GetWithPrefix(prefix string) (map[string]map[string]string, os.Error)
	Set(key string, obj map[string]string) os.Error
	Delete(key string) bool
	Update(key string, f func(map[string]string) (os.Error)) os.Error
	Inc(key, subKey string) (uint64, os.Error)
}

type resourceStorageImpl struct {
	mapStorage MapStorage
}

func NewResourceStorage(mapStorage MapStorage) ResourceStorage {
	return &resourceStorageImpl{mapStorage:mapStorage}
}

func (s *resourceStorageImpl) urlPathToStoragePath(urlPath string) string {
	if len(urlPath) <= 1 {
		return ""
	}
	if urlPath[0] != '/' {
		return ""
	}
	slashCount := strings.Count(urlPath, "/")
	if slashCount == 0 {
		return ""
	}
	return strconv.Itoa(slashCount) + ":" + urlPath
}

func (s *resourceStorageImpl) urlPathToStorageChildrenPathPrefix(urlPath string) string {
	if len(urlPath) == 0 {
		return ""
	}
	if urlPath[0] != '/' {
		return ""
	}
	if urlPath == "/" {
		return "1:/"
	}
	slashCount := strings.Count(urlPath, "/")
	if slashCount == 0 {
		return ""
	}
	return strconv.Itoa(slashCount + 1) + ":" + urlPath + "/"
}

func (s *resourceStorageImpl) Get(urlPath string) (map[string]string, os.Error) {
	// TODO: deleted attr
	storagePath := s.urlPathToStoragePath(urlPath)
	if storagePath == "" {
		return nil, nil
	}
	obj, err := s.mapStorage.Get(storagePath)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *resourceStorageImpl) GetChildren(urlPath string) (map[string]map[string]string, os.Error) {
	storagePathPrefix := s.urlPathToStorageChildrenPathPrefix(urlPath)
	if storagePathPrefix == "" {
		return nil, nil
	}
	objs := map[string]map[string]string{}
	objs2, err := s.mapStorage.GetWithPrefix(storagePathPrefix)
	if err != nil {
		return nil, err
	}
	for key, value := range objs2 {
		objs[key[len(storagePathPrefix):]] = value
	}
	return objs, nil
}

func (s *resourceStorageImpl) Set(urlPath string, obj map[string]string) os.Error {
	storagePath := s.urlPathToStoragePath(urlPath)
	if storagePath == "" {
		return nil
	}
	return s.mapStorage.Set(storagePath, obj)
}

func (s *resourceStorageImpl) Delete(urlPath string) (bool, os.Error) {
	storagePath := s.urlPathToStoragePath(urlPath)
	if storagePath == "" {
		return false, nil
	}
	/*obj, err := s.mapStorage.Get(storagePath)
	if err != nil {
		return false, err
	}
	if obj == nil {
		return false, nil
	}
	// TODO: impl*/
	/*s.mapStorage.Update(storagePath, func (obj map[string]string) os.Error {
		if obj == nil {
		}
	})*/
	return true, nil
}

func (s *resourceStorageImpl) Create(urlPath string) (uint64, os.Error) {
	storagePath := s.urlPathToStoragePath(urlPath)
	if storagePath == "" {
		return 0, nil
	}
	newID, err := s.mapStorage.Inc(urlPath, "count")
	if err != nil {
		return 0, err
	}
	return newID, nil
}
