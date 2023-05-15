package commontest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/settings"
	"mby.fr/utils/test"
)

// Common function for testing which import minimum packages to not introduce dep cycle.

const envKind = "env"
const projectKind = "project"
const imageKind = "image"

func AssertSettingsFileTree(t *testing.T, path string) {
	assert.FileExists(t, path+"/settings.yaml", "settings file should exists")
}

func AssertMinimalWorkspaceFileTree(t *testing.T, wksPath string) {
	assert.DirExists(t, wksPath, "workspace dir should exists")

	settingsDir := wksPath + "/.mass"
	assert.DirExists(t, settingsDir, ".mass/ dir should exists")
	AssertSettingsFileTree(t, settingsDir)
}

func AssertWorkspaceFileTree(t *testing.T, wksPath string) {
	AssertMinimalWorkspaceFileTree(t, wksPath)

	// Check envs FS is ok
	envsDir := wksPath + "/envs"
	assert.DirExists(t, envsDir, "config dir should exists")

	for _, env := range settings.Default().Environments {
		envPath := envsDir + "/" + env
		assert.DirExists(t, envPath, "%s config dir should exists", envPath)
		AssertEnvFileTree(t, envPath)
	}
}

func AssertEnvFileTree(t *testing.T, path string) {
	assert.DirExists(t, path, "env dir file should exists")
	assert.FileExists(t, path+"/resource.yaml", "file should exists")
	assert.FileExists(t, path+"/config.yaml", "file should exists")
}

func AssertProjectFileTree(t *testing.T, path string) {
	assert.DirExists(t, path, "project dir file should exists")
	assert.DirExists(t, path+"/test", "test dir should exists")
	//assert.FileExists(t, path + "/version.txt", "version.txt file should exists")
	assert.FileExists(t, path+"/resource.yaml", "resource file should exists")
	assert.FileExists(t, path+"/config.yaml", "config file should exists")
}

func AssertImageFileTree(t *testing.T, path string) {
	assert.DirExists(t, path, "image dir file should exists")
	assert.DirExists(t, path+"/src", "test dir should exists")
	assert.DirExists(t, path+"/test", "test dir should exists")
	//assert.FileExists(t, path+"/version.txt", "version.txt file should exists")
	assert.FileExists(t, path+"/resource.yaml", "version.txt file should exists")
	assert.FileExists(t, path+"/config.yaml", "config file should exists")
	assert.FileExists(t, path+"/Dockerfile", "buildfile file should exists")
}

func createEmptyDirs(rootDir string, dirs ...string) (err error) {
	for _, d := range dirs {
		d = filepath.Join(rootDir, d)
		//fmt.Println("Creating dir", d)
		err = os.MkdirAll(d, 0755)
		if err != nil {
			return
		}
	}
	return
}

func createEmptyFiles(rootDir string, files ...string) (err error) {
	for _, f := range files {
		parentDir := filepath.Dir(f)
		err = createEmptyDirs(rootDir, parentDir)
		if err != nil {
			return
		}
		f = filepath.Join(rootDir, f)
		//fmt.Println("Creating file", f)
		err = os.WriteFile(f, []byte(""), 0644)
		if err != nil {
			return
		}
	}
	return
}

func initConfigFile(t *testing.T, dir string) {
	f := filepath.Join(dir, "config.yaml")
	content := fmt.Sprintf("labels: \ntags: \nenvironments: \n")
	err := os.WriteFile(f, []byte(content), 0644)
	require.NoError(t, err, "Init config file should not return an error")
}

func initResourceFile(t *testing.T, dir string, kind string) {
	f := filepath.Join(dir, "resource.yaml")
	content := fmt.Sprintf("resourceKind: %s\n", kind)
	err := os.WriteFile(f, []byte(content), 0644)
	require.NoError(t, err, "Init resource file should not return an error")
}

func InitMinimalTempWorkspace(t *testing.T) (path string) {
	path, _ = test.BuildRandTempPath()
	err := settings.Init(path)
	require.NoError(t, err, "Init settings should not return an error")

	AssertMinimalWorkspaceFileTree(t, path)
	err = os.Chdir(path)
	if err != nil {
		return
	}
	return
}

func InitTempWorkspace(t *testing.T) (path string) {
	path = InitMinimalTempWorkspace(t)
	// Add envs
	for _, env := range settings.Default().Environments {
		envDir := filepath.Join(path, "envs", env)
		err := createEmptyFiles(envDir, "config.yaml", "resource.yaml")
		require.NoError(t, err, "Init temp files should not return an error")
		initConfigFile(t, envDir)
		initResourceFile(t, envDir, envKind)
	}
	err := settings.Init(path)
	require.NoError(t, err, "Init settings should not return an error")

	AssertWorkspaceFileTree(t, path)
	err = os.Chdir(path)
	if err != nil {
		return
	}
	return
}

func InitRandEnv(t *testing.T, workspacePath string) (path string) {
	name := test.RandSeq(6)
	path = filepath.Join(workspacePath, "envs", name)
	err := createEmptyFiles(path, "config.yaml", "resource.yaml")
	require.NoError(t, err, "Init temp files should not error")
	initConfigFile(t, path)
	initResourceFile(t, path, envKind)
	AssertEnvFileTree(t, path)
	return

}

func InitRandProject(t *testing.T, workspacePath string) (name, path string) {
	name = test.RandSeq(6)
	path = filepath.Join(workspacePath, name)
	err := createEmptyFiles(path, "config.yaml", "resource.yaml", "test/foo")
	require.NoError(t, err, "Init temp files should not error")
	initConfigFile(t, path)
	initResourceFile(t, path, projectKind)
	AssertProjectFileTree(t, path)
	return
}

func InitRandImage(t *testing.T, projectDir string) (name, path string) {
	name = test.RandSeq(6)
	path = filepath.Join(projectDir, name)
	err := createEmptyFiles(path, "config.yaml", "resource.yaml", "version.txt", "Dockerfile", "src/empty", "test/empty")
	require.NoError(t, err, "Init temp files should not error")
	initConfigFile(t, path)
	initResourceFile(t, path, imageKind)
	AssertImageFileTree(t, path)
	return
}
