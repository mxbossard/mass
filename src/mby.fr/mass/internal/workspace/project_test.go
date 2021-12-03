package workspace

import (
	//"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"

	"mby.fr/mass/internal/test"
)

func assertProjectFileTree(t *testing.T, path string) {
	assert.DirExists(t, path, "project dir file should exists")
	assert.DirExists(t, path + "/src", "src dir should exists")
}

func TestInitProject(t *testing.T) {
	tempDir, _ := test.BuildRandTempPath()

	os.Mkdir(tempDir, 0755)
	defer os.RemoveAll(tempDir)
	os.Chdir(tempDir)

	err := Init(tempDir)
	assert.NoError(t, err, "Init should not return an error")
	assertWorkspaceFileTree(t, tempDir)

	projectName := test.RandSeq(6)
	InitProject(projectName)
	projectPath := tempDir + "/" + projectName
	assertProjectFileTree(t, projectPath)
}

