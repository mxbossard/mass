package resource

import(
	"testing"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/workspace"
	"mby.fr/mass/internal/project"
	"mby.fr/utils/test"
)

func TestBuildResource(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	expectedName := filepath.Base(path)

	pr, err := BuildProject(path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, ProjectResource{}, pr, "bad resource type")
	assert.Equal(t, expectedName, pr.Name(), "bad resource name")
	assert.Equal(t, path, pr.Dir(), "bad resource dir")
	assert.Equal(t, ProjectKind, pr.Kind(), "bad resource kind")
}

func TestStoreResource(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	pr, err := BuildProject(path)
	require.NoError(t, err, "should not error")

	err = StoreResource(pr.BaseResource)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, defaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")
}

func TestStoreThenLoadResource(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	expectedName := filepath.Base(path)

	pr, err := BuildProject(path)
	require.NoError(t, err, "should not error")

	err = StoreResource(pr.BaseResource)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, defaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")

	loadedPr, err := LoadResource(path)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedName, loadedPr.Name(), "bad resource name")
	assert.Equal(t, path, loadedPr.Dir(), "bad resource dir")
	assert.Equal(t, ProjectKind, loadedPr.Kind(), "bad resource kind")
}

func TestInitedProjectLoadResource(t *testing.T) {
	wksPath := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(wksPath)

	prjName, prjPath := project.TestInitRandProject(t)

	expectedResourceFilepath := filepath.Join(prjPath, defaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should preexist")

	pr, err := LoadResource(prjPath)

	require.NoError(t, err, "should not error")
	assert.IsType(t, ProjectResource{}, pr, "bad resource type")
	assert.Equal(t, prjName, pr.Name(), "bad resource name")
	assert.Equal(t, prjPath, pr.Dir(), "bad resource dir")
	assert.Equal(t, ProjectKind, pr.Kind(), "bad resource kind")

}
