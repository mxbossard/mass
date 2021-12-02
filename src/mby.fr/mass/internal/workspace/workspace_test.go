package workspace

import (
	"testing"
	"os"
	"fmt"
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
	Init(wksDir)
	assertWorkspaceFileTree(t, wksPath)
}

func TestInitInExistingAbsolutePath(t *testing.T) {
	wksDir := randSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	os.Mkdir(wksPath, 0755)
	Init(wksPath)
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

	Init(wksPath)
	assertWorkspaceFileTree(t, wksPath)
}
