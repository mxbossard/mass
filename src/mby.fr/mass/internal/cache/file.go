package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"crypto/sha256"
	"encoding/hex"
	//"errors"
)

var _ = fmt.Printf

type Cache interface {
	LoadString(key string) (value string, ok bool, err error)
	StoreString(key, value string) (err error)
}

type fileCache struct {
	path string
}

func NewFileCache(path string) (cache Cache, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return
	}
	cache = fileCache{path}
	return
}

func (fc fileCache) bucketFilepath(key string) (path string) {
	hashedKey := hashKey(key)
	path = filepath.Join(fc.path, hashedKey)
	return
}

func (fc fileCache) LoadString(key string) (value string, ok bool, err error) {
	bucketPath := fc.bucketFilepath(key)
	//fmt.Printf("Loading value from bucket: %s\n", bucketPath)
	content, err := os.ReadFile(bucketPath)
	if os.IsNotExist(err) {
		err = nil
		return
	} else if err != nil {
		return
	}

	ok = true
	value = string(content)
	return
}

func (fc fileCache) StoreString(key, value string) (err error) {
	bucket := fc.bucketFilepath(key)
	//fmt.Printf("Storing value: %s in bucket: %s ...\n", value, bucket)
	err = os.WriteFile(bucket, []byte(value), 0644)
	return
}

func hashKey(key string) (h string) {
	hBytes := sha256.Sum256([]byte(key))
	h = hex.EncodeToString(hBytes[:])
	return
}

