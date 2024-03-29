package resources

import (
	//"fmt"
	"testing"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/config"
	//"mby.fr/utils/file"
	"mby.fr/utils/test"
)

func assertBaseContent(t *testing.T, path string, b Resourcer) {
	//expectedName := filepath.Base(path)
	//assert.Equal(t, expectedName, b.Name(), "bad resource name")
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

	r, err := buildEnv(path)
	require.NoError(t, err, "should not error")
	assert.NoFileExists(t, path, "should not exists")
	assertBaseContent(t, path, r.base)

	assert.Equal(t, EnvKind, r.Kind(), "bad resource kind")
}

func TestBuildProject(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	expectedName := filepath.Base(path)

	r, err := buildProject(path)
	require.NoError(t, err, "should not error")
	assert.NoFileExists(t, path, "should not exists")
	assert.Equal(t, expectedName, r.Name(), "bad resource name")
	assert.Equal(t, path, r.Dir(), "bad resource dir")
	assert.Equal(t, ProjectKind, r.Kind(), "bad resource kind")
	assertBaseContent(t, path, r)
	assertTestableContent(t, path, r)

	assert.Equal(t, ProjectKind, r.Kind(), "bad resource kind")

	_, err = r.Images()
	require.NoError(t, err, "should not error")
}

func TestBuildImage(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")

	r, err := buildImage(path)
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

/*
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
*/

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
