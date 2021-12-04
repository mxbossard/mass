package workspace

import (
	//"fmt"
	"testing"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"

	"mby.fr/utils/test"
)

func initTempWorkspace(t *testing.T) (path string) {
        path, _ = test.BuildRandTempPath()
        os.Chdir(path)
        err := Init(path)
        assert.NoError(t, err, "Init should not return an error")
        assertWorkspaceFileTree(t, path)
        return
}

func assertSettingsFileTree(t *testing.T, path string) {
	assert.FileExists(t, path + "/settings.yaml", "settings file should exists")
}

func assertConfigFileTree(t *testing.T, path string) {
}

func assertEnvFileTree(t *testing.T, path string) {
}

func assertWorkspaceFileTree(t *testing.T, wksPath string) {
	assert.DirExists(t, wksPath, "workspace dir should exists")

	settingsDir := wksPath + "/.mass"
	assert.DirExists(t, settingsDir, ".mass/ dir should exists")
	assertSettingsFileTree(t, settingsDir)

	configDir := wksPath + "/config"
	assert.DirExists(t, configDir, "config dir should exists")
	assertConfigFileTree(t, settingsDir)

	for _, env := range defaultEnvs {
		envPath := wksPath + "/config/" + env
		assert.DirExists(t, envPath, "%s config dir should exists", envPath)
		assertEnvFileTree(t, envPath)
	}
}

func TestInitInNotExistingAbsolutePath(t *testing.T) {
	wksDir := test.RandSeq(10)
	wksPath := filepath.Join(os.TempDir(), wksDir)
	os.RemoveAll(wksPath)
	defer os.RemoveAll(wksPath)

	assert.NoFileExists(t, wksPath, "workspace dir should not exists")

	err := Init(wksPath)
	assert.NoError(t, err, "Init should not return an error")
	assertWorkspaceFileTree(t, wksPath)
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
	assertWorkspaceFileTree(t, wksPath)
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
	assertWorkspaceFileTree(t, wksPath)
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
	assertWorkspaceFileTree(t, wksPath)
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
	assertWorkspaceFileTree(t, wksPath)
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
