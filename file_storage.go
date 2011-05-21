package starrpg

import (
	"os"
)

type FileStorage struct {
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path}
}

func (*FileStorage) Get(key string) ([]byte, os.Error) {
	return []byte{}, nil
}

func (*FileStorage) Set(key string) os.Error {
	return nil
}

func (*FileStorage) Delete(key string) os.Error {
	return nil
}
