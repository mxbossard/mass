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
	expectedEnvDir := filepath.Join(path, expectedEnvName)
	e, err := InitResourcer(EnvKind, expectedEnvName, path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Env{}, e, "bad type")
	assert.Equal(t, EnvKind, e.Kind(), "bad kind")
	assert.Equal(t, expectedEnvDir, e.Dir(), "bad dir")
	assert.Equal(t, expectedEnvName, e.FullName(), "bad full name")
	expectedResourceFilepath := filepath.Join(expectedEnvDir, DefaultResourceFile)
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
	r, err := Init[Env](expectedEnvName, path)
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r.directoryBase)
	expectedResourceFilepath := filepath.Join(path, expectedEnvName, DefaultResourceFile)
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
	r, err := Init[Project](expectedProjectName, path)
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

	expectedProjectName := "myProject"
	p, err := Init[Project](expectedProjectName, path)
	require.NoError(t, err, "should not error")

	expectedImageName := "myImage"
	r, err := Init[Image](expectedImageName, p)
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r)
	assertTestableFs(t, r)

	assert.DirExists(t, r.AbsSourceDir(), "source dir should exists")
	assert.FileExists(t, filepath.Join(r.Dir(), r.BuildFile), "build file should exists")
}

func TestInitProjectWithImages(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	expectedProjectName := "myProject"
	p, err := Init[Project](expectedProjectName, path)
	require.NoError(t, err, "should not error")

	assertBaseFs(t, p)
	assertTestableFs(t, p)

	images, err := p.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 0, "should not have any image")

	// Init new images
	//image1Path := filepath.Join(path, )
	_, err = Init[Image]("image1", p)
	require.NoError(t, err, "should not error")

	//image2Path := filepath.Join(path, "image2")
	_, err = Init[Image]("image2", p)
	require.NoError(t, err, "should not error")

	images, err = p.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 2, "should got 2 images")

}
