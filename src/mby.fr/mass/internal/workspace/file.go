package workspace

import (
	"os"
	"path/filepath"
)

func manageError(err error) bool {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return false
	}
	return true
}

// Create a directory. It may already exists. 
func CreateDirectory(path string) (err error) {
	err = os.MkdirAll(path, 0755)
	return
}

// Create a new directory. 
func CreateNewDirectory(path string) (err error){
	err = os.Mkdir(path, 0755)
	return
}

// Create a directory in a parent directory. 
func CreateSubDirectory(parentDirPath, name string) (path string, err error) {
	path = filepath.Join(parentDirPath, name)
	err = CreateDirectory(path)
	return
}
// Create a new directory in a parent directory. 
func CreateNewSubDirectory(parentDirPath, name string) (path string, err error) {
	path = filepath.Join(parentDirPath, name)
	err = CreateNewDirectory(path)
	return
}

// Get working directory path.
// Fail if cannot get working directory.
func WorkDirPath() (path string, err error) {
	path, err = os.Getwd()
	return
}

func Chdir(path string) (err error) {
	err = os.Chdir(path)
	return
}
