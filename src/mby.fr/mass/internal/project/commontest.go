package project

import (
	"testing"

        "github.com/stretchr/testify/assert"

        "mby.fr/utils/test"
)

func TestInitRandProject(t *testing.T) (name, path string) {
	name = test.RandSeq(6)
	path, _ = InitProject(name)
	assertProjectFileTree(t, path)
	return
}

func assertProjectFileTree(t *testing.T, path string) {
	assert.DirExists(t, path, "project dir file should exists")
	assert.DirExists(t, path + "/test", "test dir should exists")
	assert.FileExists(t, path + "/version.txt", "version.txt file should exists")
}

