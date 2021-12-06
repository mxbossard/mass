package workspace

import (
        "os"
        "testing"

        "github.com/stretchr/testify/assert"

        "mby.fr/utils/test"
)


func TestInitTempWorkspace(t *testing.T) (path string) {
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

