package resources

import (
	"fmt"
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

/*
func TestWriteDirectoryBase(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	b, err := buildDirectoryBase(EnvKind, path)
	require.NoError(t, err, "should not error")

	err = Write(b)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := b.backingFilepath()
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

	te, err := buildProject(path)
	require.NoError(t, err, "should not error")

	err = Write(te)
	require.NoError(t, err, "should not error")

	expectedResourceFilepath := filepath.Join(path, DefaultResourceFile)
	assert.FileExists(t, expectedResourceFilepath, "resource file should exist")
}
*/
func TestWriteProject(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	pr, err := buildProject(path)
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

	content := "resourceKind: env\n"
	resourceFilepath := filepath.Join(path, DefaultResourceFile)
	err = os.WriteFile(resourceFilepath, []byte(content), 0644)
	require.NoError(t, err, "should not error")

	r, err := ReadAny(path)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Env{}, r, "bad type")
}

func TestReadResourcer(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	content := `
        resourceKind: project
    `
	expectedResName := "bar"
	expectedResDir := filepath.Join(path, expectedResName)
	err = os.MkdirAll(expectedResDir, 0755)
	require.NoError(t, err, "should not error")
	resourceFilepath := filepath.Join(expectedResDir, DefaultResourceFile)
	err = os.WriteFile(resourceFilepath, []byte(content), 0644)
	require.NoError(t, err, "should not error")

	r, err := ReadResourcer(resourceFilepath)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Project{}, r, "bad type")

	assert.Equal(t, expectedResDir, r.Dir())
	assert.Equal(t, expectedResName, r.FullName())

	if p, ok := r.(Project); ok {
		assert.Equal(t, expectedResName, p.FullName(), "bad name")
		testFunc := func() {
			p.Images()
		}
		assert.NotPanics(t, testFunc, "should panic")

		assert.Equal(t, expectedResName, (&p).FullName())
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

	expectedProjectName := "projectName"
	expectedProjectDir := filepath.Join(path, expectedProjectName)
	os.MkdirAll(expectedProjectDir, 0755)
	content := "resourceKind: project\n"
	resourceFilepath := filepath.Join(expectedProjectDir, DefaultResourceFile)
	err = os.WriteFile(resourceFilepath, []byte(content), 0644)
	require.NoError(t, err, "should not error")

	_, err = Read[Env](expectedProjectDir)
	assert.Error(t, err, "should error")

	p, err := Read[Project](expectedProjectDir)
	require.NoError(t, err, "should not error")
	assert.IsType(t, Project{}, p, "bad type")

	assert.Equal(t, expectedProjectDir, p.Dir(), "bad dir")
	assert.Equal(t, expectedProjectName, p.FullName(), "bad name")
	testFunc := func() {
		p.Images()
	}
	assert.NotPanics(t, testFunc, "should panic")
}

func TestWriteThenRead(t *testing.T) {
	path, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	expectedImageName := "monImage"
	expectedProjectName := "monProjet"
	expectedImageFullName := fmt.Sprintf("%s/%s", expectedProjectName, expectedImageName)
	expectedProjectDir := filepath.Join(path, expectedProjectName)
	expectedImageDir := filepath.Join(expectedProjectDir, expectedImageName)
	p, err := Build[Project](expectedProjectName, path)
	require.NoError(t, err, "should not error")
	i, err := Build[Image](expectedImageName, p)
	require.NoError(t, err, "should not error")
	//i.Init()

	err = Write(p)
	require.NoError(t, err, "should not error")

	err = Write(i)
	require.NoError(t, err, "should not error")

	expectedImageResFilepath := filepath.Join(expectedImageDir, DefaultResourceFile)
	assert.FileExists(t, expectedImageResFilepath, "resource file should exist")

	loadedImage, err := Read[Image](expectedImageResFilepath)
	require.NoError(t, err, "should not error")
	//loadedImage := res.(*Image)
	assert.Equal(t, expectedImageDir, loadedImage.Dir(), "bad resource dir")
	assert.Equal(t, ImageKind, loadedImage.Kind(), "bad resource kind")
	assert.Equal(t, filepath.Join(expectedImageDir, DefaultSourceDir), loadedImage.AbsSourceDir(), "bad source dir")
	assert.Equal(t, DefaultBuildFile, loadedImage.BuildFile, "bad build file")
	assert.Equal(t, filepath.Join(expectedImageDir, DefaultBuildFile), loadedImage.AbsBuildFile(), "bad build file")
	assert.Equal(t, DefaultInitialVersion, loadedImage.Version(), "bad version")
	assert.Equal(t, expectedImageFullName, loadedImage.FullName(), "bad image full name")

	project, err := loadedImage.Project()
	require.NoError(t, err, "should not error")

	assert.NotNil(t, project, "bad parent project")
	assert.Equal(t, ProjectKind, project.Kind(), "bad parent project kind")
	assert.Equal(t, expectedProjectDir, project.Dir(), "bad parent project dir")
	assert.Equal(t, expectedProjectName, project.FullName(), "bad project full name")
}
