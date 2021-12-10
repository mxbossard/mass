package project

import (
	"testing"

        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"

        "mby.fr/utils/test"
)

func TestInitRandProject(t *testing.T) (name, path string) {
	name = test.RandSeq(6)
	path, err := InitProject(name)
	require.NoError(t, err, "should not error")
	assertProjectFileTree(t, path)
	return
}

func assertProjectFileTree(t *testing.T, path string) {
	assert.DirExists(t, path, "project dir file should exists")
	assert.DirExists(t, path + "/test", "test dir should exists")
	assert.FileExists(t, path + "/version.txt", "version.txt file should exists")
	assert.FileExists(t, path + "/resource.yaml", "resource file should exists")
	assert.FileExists(t, path + "/config.yaml", "config file should exists")
}

func TestInitRandImage(t *testing.T, p Project) (name, path string) {
	name = test.RandSeq(6)
	path, err := InitImage(p, name)
	require.NoError(t, err, "should not error")
	assertImageFileTree(t, path)
	return
}

func assertImageFileTree(t *testing.T, path string) {
	assert.DirExists(t, path, "image dir file should exists")
	assert.DirExists(t, path + "/src", "test dir should exists")
	assert.DirExists(t, path + "/test", "test dir should exists")
	assert.FileExists(t, path + "/version.txt", "version.txt file should exists")
	assert.FileExists(t, path + "/resource.yaml", "version.txt file should exists")
	assert.FileExists(t, path + "/config.yaml", "config file should exists")
}

