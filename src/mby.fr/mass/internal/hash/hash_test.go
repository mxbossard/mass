package hash

import (
	//"fmt"
	"testing"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"

	"mby.fr/mass/internal/test"
	"mby.fr/mass/internal/workspace"
)

func assertModifiedDir(t *testing.T, path, message string) {
	ok, err := IsDirectoryModified(path)
	assert.NoError(t, err, "should not produce an error: %s", message)
	assert.True(t, ok, "dir1 should be marked as modified: %s", message)
}

func assertNotModifiedDir(t *testing.T, path, message string) {
	ok, err := IsDirectoryModified(path)
	assert.NoErrorf(t, err, "should not produce an error: %s", message)
	assert.Falsef(t, ok, "dir1 should be marked as modified: %s", message)
}

func InitTempWorspace() (path string) {
	path, _ = test.BuildRandTempPath()
	workspace.Init(path)
	os.Chdir(path)
	return
}

func TestHashEmptyDir(t *testing.T) {
	t.Skip()
	path := InitTempWorspace()
	defer os.RemoveAll(path)
	assert.DirExists(t, path, "Temp workspace dir should exists")

	dir1 := filepath.Join(path, "dir1")
	os.Mkdir(dir1, 0755)
	assertModifiedDir(t, dir1, "empty dir1")

	assertNotModifiedDir(t, dir1, "empty dir1")
}

func TestHashDir(t *testing.T) {
	t.Skip()
	path := InitTempWorspace()
	defer os.RemoveAll(path)
	assert.DirExists(t, path, "Temp workspace dir should exists")

	dir1 := filepath.Join(path, "dir1")
	os.Mkdir(dir1, 0755)
	assertModifiedDir(t, dir1, "empty dir1")

	assertNotModifiedDir(t, dir1, "empty dir1")

	file1 := filepath.Join(dir1, "file1")
	os.WriteFile(file1, []byte(""), 0644)
	assertModifiedDir(t, dir1, "empty file1")

	assertNotModifiedDir(t, dir1, "empty file1")

	file2 := filepath.Join(dir1, "file2")
	os.WriteFile(file2, []byte("foo"), 0644)
	assertModifiedDir(t, dir1, "foo file2")

	assertNotModifiedDir(t, dir1, "foo file2")

	os.WriteFile(file2, []byte("bar"), 0644)
	assertModifiedDir(t, dir1, "bar file2")

	assertNotModifiedDir(t, dir1, "bar file2")
}

