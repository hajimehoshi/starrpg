package starrpg

import (
	"io/ioutil"
	"os"
	"sync"
)

type fileCacheEntry struct {
	content []byte
	mtime_ns int64
}

type fileCache struct {
	entries map[string]*fileCacheEntry
	lock *sync.RWMutex
}

var (
	fileCache_ = fileCache{
	entries: map[string]*fileCacheEntry{},
	lock: &sync.RWMutex{},}
)

func (c *fileCache) getEntry(path string) (*fileCacheEntry, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	ent, ok := c.entries[path]
	return ent, ok
}

func (c *fileCache) insertEntry(path string, ent *fileCacheEntry) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.entries[path] = ent
}

func (c *fileCache) deleteEntry(path string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.entries[path] = nil, false
}

func GetFileFromCache(filepath string) ([]byte, os.Error) {
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		fileCache_.deleteEntry(filepath)
		if e, ok := err.(*os.PathError); ok && e.Error == os.ENOENT {
			return nil, nil
		}
		return nil, err
	}
	if !fileinfo.IsRegular() {
		fileCache_.deleteEntry(filepath)
		return nil, nil
	}
	ent, hit := fileCache_.getEntry(filepath)
	if !hit || fileinfo.Mtime_ns != ent.mtime_ns {
		fileCache_.deleteEntry(filepath)
		file, err := os.OpenFile(filepath, os.O_RDONLY, 0777)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		ent = &fileCacheEntry{content, fileinfo.Mtime_ns}
		fileCache_.insertEntry(filepath, ent)
	}
	return ent.content, nil
}
