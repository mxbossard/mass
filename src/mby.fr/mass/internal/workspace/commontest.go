package workspace

import (
        "os"
        "testing"

        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"

        "mby.fr/utils/test"
        "mby.fr/mass/internal/settings"
)


func TestInitTempWorkspace(t *testing.T) (path string) {
        path, _ = test.BuildRandTempPath()
        os.Chdir(path)
        err := Init(path)
        require.NoError(t, err, "Init should not return an error")
	assertWorkspaceFileTree(t, path)
        return
}

func assertSettingsFileTree(t *testing.T, path string) {
	assert.FileExists(t, path + "/settings.yaml", "settings file should exists")
}

func assertWorkspaceFileTree(t *testing.T, wksPath string) {
	assert.DirExists(t, wksPath, "workspace dir should exists")

	settingsDir := wksPath + "/.mass"
	assert.DirExists(t, settingsDir, ".mass/ dir should exists")
	assertSettingsFileTree(t, settingsDir)

	configDir := wksPath + "/config"
	assert.DirExists(t, configDir, "config dir should exists")
	//assertConfigFileTree(t, settingsDir)

	for _, env := range settings.Default().Environments {
		envPath := wksPath + "/config/" + env
		assert.DirExists(t, envPath, "%s config dir should exists", envPath)
		assertEnvFileTree(t, envPath)
	}
}

func TestInitRandEnv(t *testing.T) (path string) {
	name := test.RandSeq(6)
        path, err := InitEnv(name)
        require.NoError(t, err, "should not error")
        assertEnvFileTree(t, path)
        return

}

func assertEnvFileTree(t *testing.T, path string) {
        assert.DirExists(t, path, "env dir file should exists")
        assert.FileExists(t, path + "/resource.yaml", "file should exists")
        assert.FileExists(t, path + "/config.yaml", "file should exists")
}

