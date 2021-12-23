package workspace

import (
	//"fmt"
	"testing"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"

	"mby.fr/utils/test"
	"mby.fr/mass/internal/commontest"
)

func TestInitInNotExistingAbsolutePath(t *testing.T) {
	wksDir := test.RandSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	err := Init(wksPath)
	assert.NoError(t, err, "Init should not return an error")
	commontest.AssertWorkspaceFileTree(t, wksPath)
}

func TestInitInNotExistingRelativePath(t *testing.T) {
	wksDir := test.RandSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	os.Chdir(os.TempDir())
	err := Init(wksDir)
	assert.NoError(t, err, "Init should not return an error")
	commontest.AssertWorkspaceFileTree(t, wksPath)
}

func TestInitInExistingAbsolutePath(t *testing.T) {
	wksDir := test.RandSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	os.Mkdir(wksPath, 0755)
	err := Init(wksPath)
	assert.NoError(t, err, "Init should not return an error")
	commontest.AssertWorkspaceFileTree(t, wksPath)
}

func TestInitInNotExistingAbsoluteSubPath(t *testing.T) {
	wksDir := test.RandSeq(10)
	parentDir := test.RandSeq(10)
	parentPath := filepath.Join(os.TempDir(), parentDir)
	wksPath := filepath.Join(parentPath, wksDir)
	os.RemoveAll(parentPath)
	defer os.RemoveAll(parentPath)

	assert.NoFileExists(t, parentPath, "parent dir should not exists")
	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	err := Init(wksPath)
	assert.NoError(t, err, "Init should not return an error")
	commontest.AssertWorkspaceFileTree(t, wksPath)
}

func TestInitWithDotPath(t *testing.T) {
	wksDir := test.RandSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")
	os.Mkdir(wksPath, 0755)
	os.Chdir(wksPath)

	dotPath := "."
	err := Init(dotPath)
	assert.NoError(t, err, "Init should not return an error")
	commontest.AssertWorkspaceFileTree(t, wksPath)
}

func TestInitWithEmptyPath(t *testing.T) {
	wksDir := test.RandSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")
	os.Mkdir(wksPath, 0755)
	os.Chdir(wksPath)

	emptyPath := ""
	err := Init(emptyPath)
	assert.Error(t, err, "Init should return an error")
	assert.NoFileExists(t, wksPath, "workspace dir should not exists")
}
