package resources

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/test"
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

func assertBaseContent(t *testing.T, path string, b Resourcer) {
	expectedName := filepath.Base(path)
	assert.Equal(t, expectedName, b.Name(), "bad resource name")
	assert.Equal(t, path, b.Dir(), "bad resource dir")
}

func assertBaseFs(t *testing.T, b Resourcer) {
	assert.DirExists(t, b.Dir(), "should exists")
	assert.FileExists(t, b.Dir()+"/"+config.DefaultConfigFile, "should exists")
	assert.FileExists(t, b.Dir()+"/"+DefaultResourceFile, "should exists")
}

func assertTestableContent(t *testing.T, path string, r Resourcer) {
	//tester, _ := Undecorate(r, testable{})
	//require.Implements(t, (*Tester)(nil), r, "Should implements Tester")
	tester, ok := r.(Tester)
	require.True(t, ok, "Should implements Tester !")
	assert.Equal(t, path+"/"+DefaultTestDir, tester.AbsTestDir(), "bad resource dir")
}

func assertTestableFs(t *testing.T, r Resourcer) {
	//tester, _ := Undecorate(r, testable{})
	//require.Implements(t, (*Tester)(nil), r, "Should implements Tester")
	tester, ok := r.(Tester)
	require.True(t, ok, "Should implements Tester !")
	assert.DirExists(t, tester.AbsTestDir(), "should exists")
}

func TestBuildEnv(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildEnv(path)
	require.NoError(t, err, "should not error")
	assert.NoFileExists(t, path, "should not exists")
	assertBaseContent(t, path, r.base)

	assert.Equal(t, EnvKind, r.Kind(), "bad resource kind")
}

func TestInitEnv(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildEnv(path)
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	err = r.Init()
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r.base)
}

func TestBuildProject(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildProject(path)
	require.NoError(t, err, "should not error")
	assert.NoFileExists(t, path, "should not exists")
	assertBaseContent(t, path, r)
	assertTestableContent(t, path, r)

	assert.Equal(t, ProjectKind, r.Kind(), "bad resource kind")

	_, err = r.Images()
	require.NoError(t, err, "should not error")
}

func TestInitProject(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	//r, err := Init2[Project](ProjectKind, path)
	r, err := FromPath[Project](path)
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	err = r.Init()
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r)
	assertTestableFs(t, r)

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

	assertBaseContent(t, path, r)
	assertTestableContent(t, path, r)

	assert.Equal(t, ImageKind, r.Kind(), "bad resource kind")
	assert.Equal(t, path+"/"+DefaultSourceDir, r.AbsSourceDir(), "bad source dir")
	assert.Equal(t, path+"/"+DefaultBuildFile, r.AbsBuildFile(), "bad buildfile")
	assert.Equal(t, DefaultInitialVersion, r.Version(), "bad version")

	parentDir := filepath.Dir(path)
	assert.NotNil(t, r.Project, "bad parent project")
	assert.Equal(t, ProjectKind, r.Project.Kind(), "bad parent project kind")
	assert.Equal(t, parentDir, r.Project.Dir(), "bad parent project dir")
}

func TestInitImage(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildImage(path)
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	err = r.Init()
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r)
	assertTestableFs(t, r)

	assert.DirExists(t, r.AbsSourceDir(), "source dir should exists")
	assert.FileExists(t, r.BuildFile, "source dir should exists")
}

func TestInitProjectWithImages(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := BuildProject(path)
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	err = r.Init()
	require.NoError(t, err, "should not error")

	assertBaseFs(t, r)
	assertTestableFs(t, r)

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
