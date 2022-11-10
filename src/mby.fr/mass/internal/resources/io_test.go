package resources

import (
	//"fmt"
	"os"
	"testing"

	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/utils/test"
)

func assertFileContains(t *testing.T, path string, expectedContent string) {
	path, err := filepath.Abs(path)
	require.NoError(t, err, "should not error")
	content, err := os.ReadFile(path)
	require.NoError(t, err, "should not error")
	stringContent := string(content)
	assert.Equal(t, expectedContent, stringContent, "Bad file content")
}

func TestWriteBase(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	b, err := buildBase(EnvKind, path)
	require.NoError(t, err, "should not error")

	err = Write(b)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, DefaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")
	//base := filepath.Base(path)
	//expectedContent := fmt.Sprintf("{1 %s %s}\n", base, path)
	expectedContent := "resourceKind: env\n"
	assertFileContains(t, expectedResourceFilepath, expectedContent)
}

func TestWriteTestable(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	te, err := BuildProject(path)
	require.NoError(t, err, "should not error")

	err = Write(te)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, DefaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")
}

func TestWriteProject(t *testing.T) {
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

func TestReadAny(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	envContent := "resourceKind: env\n"
	err = os.WriteFile(path, []byte(envContent), 0644)

	r, err := ReadAny(path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Env{}, r, "bad type")
}

func TestReadResourcer(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	envContent := "resourceKind: project\n"
	err = os.WriteFile(path, []byte(envContent), 0644)

	r, err := ReadResourcer(path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Project{}, r, "bad type")

	expectedName := filepath.Base(path)
	assert.Equal(t, expectedName, r.Name())

	if p, ok := r.(Project); ok {
		expectedName := filepath.Base(path)
		assert.Equal(t, expectedName, p.Name(), "bad name")
		testFunc := func() {
			p.Images()
		}
		assert.Panics(t, testFunc, "should panic")

		assert.Equal(t, expectedName, (&p).Name())
		testFunc = func() {
			(&p).Images()
		}
		assert.NotPanics(t, testFunc, "should not panic")
	}
}

func TestRead(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	envContent := "resourceKind: project\n"
	err = os.WriteFile(path, []byte(envContent), 0644)

	_, err = Read[Env](path)
	assert.Error(t, err, "should error")

	p, err := Read[Project](path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Project{}, p, "bad type")

	expectedName := filepath.Base(path)
	assert.Equal(t, expectedName, p.Name(), "bad name")
	testFunc := func() {
		p.Images()
	}
	assert.Panics(t, testFunc, "should panic")

	ePtr, err := Read[*Project](path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Project{}, ePtr, "bad type")

	assert.Equal(t, expectedName, ePtr.Name())
	assert.NotPanics(t, testFunc, "should not panic")
}

func TestWriteThenRead(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	i, err := BuildImage(path)
	require.NoError(t, err, "should not error")
	//i.Init()
	err = Write(i)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, DefaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")

	loadedImage, err := Read[Image](path)
	require.NoError(t, err, "should not error")
	//loadedImage := res.(*Image)
	assert.Equal(t, path, loadedImage.Dir(), "bad resource dir")
	assert.Equal(t, ImageKind, loadedImage.Kind(), "bad resource kind")
	assert.Equal(t, path+"/"+DefaultSourceDir, loadedImage.AbsSourceDir(), "bad source dir")
	assert.Equal(t, DefaultBuildFile, loadedImage.BuildFile, "bad build file")
	assert.Equal(t, path+"/"+DefaultBuildFile, loadedImage.AbsBuildFile(), "bad build file")
	assert.Equal(t, DefaultInitialVersion, loadedImage.Version(), "bad version")

	parentDir := filepath.Dir(path)
	assert.NotNil(t, loadedImage.Project, "bad parent project")
	assert.Equal(t, ProjectKind, loadedImage.Project.Kind(), "bad parent project kind")
	assert.Equal(t, parentDir, loadedImage.Project.Dir(), "bad parent project dir")
}
