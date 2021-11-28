package workspace

import (
	"log"
	"os"
	"path/filepath"
)

// Create a new directory. 
// Fail if cannot create directory or directory already exists.
func CreateNewDirectory(parentDir, name string) {
	newDir := filepath.Join(parentDir, name)

	err := os.Mkdir(newDir, 0755)
        if (err != nil) {
                log.Fatal(err)
        }
}


