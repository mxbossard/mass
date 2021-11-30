package workspace

import (
	"log"
	"os"
	"path/filepath"
)

// Create a new directory. 
// Fail if cannot create directory or directory already exists.
func CreateNewDirectory(path string) {
	err := os.Mkdir(path, 0755)
        if (err != nil) {
                log.Fatal(err)
        }
}

// Create a new directory in a parent directory. 
// Fail if cannot create directory or directory already exists.
func CreateNewSubDirectory(parentDirPath, name string) {
	newDirPath := filepath.Join(parentDirPath, name)
	CreateNewDirectory(newDirPath)
}

// Get working directory path.
// Fail if cannot get working directory.
func GetWorkDirPath() string {
        workDirPath, err := os.Getwd()
        if err != nil {
                log.Fatal(err)
        }
        return workDirPath
}

func Chdir(path string) {
	err := os.Chdir(path)
        if (err != nil) {
                log.Fatal(err)
        }
}
