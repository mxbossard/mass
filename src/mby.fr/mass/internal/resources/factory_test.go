package resources

import (
	//"os"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	//"mby.fr/mass/internal/settings"
	//"mby.fr/utils/file"
	"mby.fr/utils/test"
)

func TestBuildAny(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	e, err := BuildAny(EnvKind, path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Env{}, e, "bad type")
}

func TestBuildResourcer(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	e, err := BuildResourcer(EnvKind, path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Env{}, e, "bad type")
}

func TestBuild(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	e, err := Build[Env](path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Env{}, e, "bad type")
	assert.Equal(t, EnvKind, e.Kind(), "bad kind")
	assert.Equal(t, path, e.Dir(), "bad dir")
}

func TestCallFuncOnResource(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	e, err := Build[Env](path)
	call := func(r Resourcer) (res Resourcer, err error) {
		fmt.Printf("call on type: %T", r)
		res = r
		return
	}
	r, err := CallFuncOnResource[Env](e, call)
	require.NoError(t, err, "should not error")
	assert.Equal(t, e.Name(), r.Name())
}

/*
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

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

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
	require.NoError(t, err, "should not error")
}

func TestInitProject(t *testing.T) {
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

	assertBaseFs(t, r.Base)
	assertTestableFs(t, r.Testable)

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
*/
