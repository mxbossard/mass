package resources

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/settings"
	"mby.fr/utils/test"
)

func TestInitResourcer(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	expectedEnvName := "myEnv"
	e, err := InitResourcer(EnvKind, path, expectedEnvName)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Env{}, e, "bad type")
	assert.Equal(t, EnvKind, e.Kind(), "bad kind")
	assert.Equal(t, path, e.Dir(), "bad dir")
	assert.Equal(t, expectedEnvName, e.FullName(), "bad full name")
	expectedResourceFilepath := filepath.Join(path, DefaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "Resource file not written")
}

func TestInitEnv(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	expectedEnvName := "myEnv"
	r, err := Init[Env](path, expectedEnvName)
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r.directoryBase)
	expectedResourceFilepath := filepath.Join(path, DefaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "Resource file not written")
}

func TestInitProject(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	expectedProjectName := "myProject"
	r, err := Init[Project](path, expectedProjectName)
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r)
	assertTestableFs(t, r)

	images, err := r.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 0, "should not have any image")
}

func TestInitImage(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	expectedImageName := "myImage"
	r, err := Init[Image](path, expectedImageName)
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r)
	assertTestableFs(t, r)

	assert.DirExists(t, r.AbsSourceDir(), "source dir should exists")
	assert.FileExists(t, r.BuildFile, "source dir should exists")
}

func TestInitProjectWithImages(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	expectedProjectName := "myProject"
	p, err := Init[Project](path, expectedProjectName)
	require.NoError(t, err, "should not error")

	assertBaseFs(t, p)
	assertTestableFs(t, p)

	images, err := p.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 0, "should not have any image")

	// Init new images
	//image1Path := filepath.Join(path, )
	_, err = Init[Image](p.Dir(), "image1")
	require.NoError(t, err, "should not error")

	//image2Path := filepath.Join(path, "image2")
	_, err = Init[Image](p.Dir(), "image2")
	require.NoError(t, err, "should not error")

	images, err = p.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 2, "should got 2 images")

}
