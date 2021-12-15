package resources

import(
	"testing"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/utils/test"
	"mby.fr/mass/internal/config"
)

func TestBuildResources(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	expectedName := filepath.Base(path)

	pr, err := BuildProject(path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Project{}, pr, "bad resource type")
	assert.Equal(t, expectedName, pr.Name(), "bad resource name")
	assert.Equal(t, path, pr.Dir(), "bad resource dir")
	assert.Equal(t, ProjectKind, pr.Kind(), "bad resource kind")
}

func TestWrite(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	pr, err := BuildProject(path)
	require.NoError(t, err, "should not error")

	err = Write(pr)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, DefaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")
}

func TestWriteThenRead(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	expectedName := filepath.Base(path)

	i, err := BuildImage(path)
	require.NoError(t, err, "should not error")

	err = Write(i)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, DefaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")

	res, err := Read(path)
	require.NoError(t, err, "should not error")
	loadedImage := res.(Image)
	assert.Equal(t, expectedName, loadedImage.Name(), "bad resource name")
	assert.Equal(t, path, loadedImage.Dir(), "bad resource dir")
	assert.Equal(t, ImageKind, loadedImage.Kind(), "bad resource kind")
	assert.Equal(t, path + "/" + DefaultSourceDir, loadedImage.SourceDir(), "bad source dir")
	assert.Equal(t, path + "/" + DefaultBuildFile, loadedImage.BuildFile, "bad build file")
	assert.Equal(t, DefaultInitialVersion, loadedImage.Version, "bad version")
}

func assertBaseContent(t *testing.T, path string, b Base) {
	expectedName := filepath.Base(path)
	assert.Equal(t, expectedName, b.Name(), "bad resource name")
	assert.Equal(t, path, b.Dir(), "bad resource dir")
}

func assertBaseFs(t *testing.T, b Base) {
	assert.DirExists(t, b.Dir(), "should exists")
	assert.FileExists(t, b.Dir() + "/" + config.DefaultConfigFile, "should exists")
	assert.FileExists(t, b.Dir() + "/" + DefaultResourceFile, "should exists")
}

func assertTestableContent(t *testing.T, path string, r Testable) {
	assert.Equal(t, path + "/" + DefaultTestDir, r.TestDir(), "bad resource dir")
}

func assertTestableFs(t *testing.T, r Testable) {
	assert.DirExists(t, r.TestDir(), "should exists")
}

func TestBuildEnv(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildEnv(path)
	require.NoError(t, err, "should not error")
	assert.NoFileExists(t, path, "should not exists")
	assertBaseContent(t, path, r.Base)

	assert.Equal(t, EnvKind, r.Kind(), "bad resource kind")
}

func TestInitEnv(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildEnv(path)
	require.NoError(t, err, "should not error")

	err = r.Init()
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r.Base)
}

func TestBuildProject(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildProject(path)
	require.NoError(t, err, "should not error")
	assert.NoFileExists(t, path, "should not exists")
	assertBaseContent(t, path, r.Base)
	assertTestableContent(t, path, r.Testable)

	assert.Equal(t, ProjectKind, r.Kind(), "bad resource kind")

	_, err = r.Images()
	require.Error(t, err, "should error")
}

func TestInitProject(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildProject(path)
	require.NoError(t, err, "should not error")

	err = r.Init()
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r.Base)
	assertTestableFs(t, r.Testable)

	images, err := r.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 0, "should not have any image")
}

func TestBuildImage(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildImage(path)
	require.NoError(t, err, "should not error")
	assert.NoFileExists(t, path, "should not exists")

	assertBaseContent(t, path, r.Base)
	assertTestableContent(t, path, r.Testable)

	assert.Equal(t, ImageKind, r.Kind(), "bad resource kind")
	assert.Equal(t, path + "/" + DefaultSourceDir, r.SourceDir(), "bad source dir")
	assert.Equal(t, path + "/" + DefaultBuildFile, r.BuildFile, "bad buildfile")
	assert.Equal(t, DefaultInitialVersion, r.Version, "bad version")
}

func TestInitImage(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildImage(path)
	require.NoError(t, err, "should not error")

	err = r.Init()
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r.Base)
	assertTestableFs(t, r.Testable)

	assert.DirExists(t, r.SourceDir(), "source dir should exists")
	assert.FileExists(t, r.BuildFile, "source dir should exists")
}

func TestInitProjectWithImages(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildProject(path)
	require.NoError(t, err, "should not error")

	err = r.Init()
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r.Base)
	assertTestableFs(t, r.Testable)

	images, err := r.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 0, "should not have any image")

	// Init new images
	image1Path := filepath.Join(path, "image1")
	i1, err := BuildImage(image1Path)
	i1.Init()

	image2Path := filepath.Join(path, "image2")
	i2, err := BuildImage(image2Path)
	i2.Init()

	images, err = r.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 2, "should got 2 images")

}

