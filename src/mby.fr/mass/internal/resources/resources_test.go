package resources

import(
	"testing"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/utils/test"
)

func TestBuildResources(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	expectedName := filepath.Base(path)

	pr, err := buildProject(path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Project{}, pr, "bad resource type")
	assert.Equal(t, expectedName, pr.Name(), "bad resource name")
	assert.Equal(t, path, pr.Dir(), "bad resource dir")
	assert.Equal(t, ProjectKind, pr.Kind(), "bad resource kind")
}

func TestStore(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	pr, err := buildProject(path)
	require.NoError(t, err, "should not error")

	err = Store(pr.Base)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, defaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")
}

func TestStoreThenLoad(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	expectedName := filepath.Base(path)

	pr, err := buildProject(path)
	require.NoError(t, err, "should not error")

	err = Store(pr.Base)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, defaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")

	loadedPr, err := Load(path)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedName, loadedPr.Name(), "bad resource name")
	assert.Equal(t, path, loadedPr.Dir(), "bad resource dir")
	assert.Equal(t, ProjectKind, loadedPr.Kind(), "bad resource kind")
}

