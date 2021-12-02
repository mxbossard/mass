package workspace

import (
	"fmt"
	"testing"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"time"
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func init() {
	fmt.Println("Rand seed initialization ...")
	rand.Seed(time.Now().UnixNano())
}


func assertWorkspaceFileTree(t *testing.T, wksPath string) {
	assert.DirExists(t, wksPath, "workspace dir should exists")
	assert.DirExists(t, wksPath + "/.mass", ".mass/ dir should exists")
	assert.FileExists(t, wksPath + "/.mass/settings.yaml", "settings file should exists")
	assert.DirExists(t, wksPath + "/config", "config dir should exists")
	for _, env := range defaultEnvs {
		envPath := wksPath + "/config/" + env
		assert.DirExists(t, envPath, "%s config dir should exists", envPath)
	}
}

func TestInitInNotExistingAbsolutePath(t *testing.T) {
	wksDir := randSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	err := Init(wksPath)
	assert.NoError(t, err, "Init should not return an error")
	assertWorkspaceFileTree(t, wksPath)
}

func TestInitInNotExistingRelativePath(t *testing.T) {
	wksDir := randSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	os.Chdir(os.TempDir())
	err := Init(wksDir)
	assert.NoError(t, err, "Init should not return an error")
	assertWorkspaceFileTree(t, wksPath)
}

func TestInitInExistingAbsolutePath(t *testing.T) {
	wksDir := randSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	os.Mkdir(wksPath, 0755)
	err := Init(wksPath)
	assert.NoError(t, err, "Init should not return an error")
	assertWorkspaceFileTree(t, wksPath)
}

func TestInitInNotExistingAbsoluteSubPath(t *testing.T) {
	wksDir := randSeq(10)
	parentDir := randSeq(10)
	parentPath := filepath.Join(os.TempDir(), parentDir)
	wksPath := filepath.Join(parentPath, wksDir)
	os.RemoveAll(parentPath)
	defer os.RemoveAll(parentPath)

	assert.NoFileExists(t, parentPath, "parent dir should not exists")
	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	err := Init(wksPath)
	assert.NoError(t, err, "Init should not return an error")
	assertWorkspaceFileTree(t, wksPath)
}

func TestInitWithDotPath(t *testing.T) {
	wksDir := randSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")
	os.Mkdir(wksPath, 0755)
	os.Chdir(wksPath)

	dotPath := "."
	err := Init(dotPath)
	assert.NoError(t, err, "Init should not return an error")
	assertWorkspaceFileTree(t, wksPath)
}

func TestInitWithEmptyPath(t *testing.T) {
	wksDir := randSeq(10)
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
