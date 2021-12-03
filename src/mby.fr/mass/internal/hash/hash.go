package hash

import (
	"fmt"
	"os"
	"path/filepath"
	"crypto/sha256"
	"encoding/hex"
	//"errors"

	"golang.org/x/mod/sumdb/dirhash"

	"mby.fr/mass/internal/workspace"
)

func IsDirectoryModified(path string) (ok bool, err error) {
	h, err := dirhash.HashDir(path, "", dirhash.Hash1)
	fmt.Printf("path: %s => currentHash: %s)\n", path, h)
	if err != nil {
		return
	}

	previousHash, err := LoadPreviousHash(path)

	err2 := StoreHash(path, h)
	if err2 != nil {
		return false, err2
	}

	fmt.Printf("path: %s => previousHash: %s)\n", path, previousHash)
	if err != nil || previousHash == "" {
		ok = true
		return
	}

	ok = h == previousHash
	return
}

func HashPath(path string) (h string) {
	hBytes := sha256.Sum256([]byte(path))
	h = hex.EncodeToString(hBytes[:])
	return
}

func HashStorePath(path string) (hStorePath string, err error) {
	ss, err := workspace.GetSettingsService()
	hashDir := ss.HashDirPath()

	hPath := HashPath(path)
	hStorePath = filepath.Join(hashDir, hPath)

	return
}

func LoadPreviousHash(path string) (h string, err error) {
	hStorePath, err := HashStorePath(path)
	if err != nil {
		return
	}
	fmt.Printf("Loading previous hash from: %s ...\n", hStorePath)
	content, err := os.ReadFile(hStorePath)  
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return
	}

	h = hex.EncodeToString(content)
	//h = string(content)
	return
}

func StoreHash(path, h string) (err error) {
	hStorePath, err := HashStorePath(path)
	if err != nil {
		return
	}
	fmt.Printf("Storing hash: %s in: %s ...\n", h, hStorePath)
	err = os.WriteFile(hStorePath, []byte(h), 0644)
	return
}
