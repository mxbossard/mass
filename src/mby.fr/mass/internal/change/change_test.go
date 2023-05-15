package change

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/test"
)

func TestInit(t *testing.T) {
	tempDir, err := test.MkRandTempDir()
	defer os.RemoveAll(tempDir)
	require.NoError(t, err, "should not error")
	require.NoFileExists(t, tempDir, "should not exists")

	// Init Settings for templates to work
	err = settings.Init(tempDir)
	require.NoError(t, err, "should not error")
	os.Chdir(tempDir)

	err = Init()
	require.NoError(t, err, "should not error")
	assert.DirExists(t, tempDir, "should exists")
}

func TestCalcImageSignature(t *testing.T) {
	path, err := test.BuildRandTempPath()
	defer os.RemoveAll(path)
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	expectedImageName := "monImage"
	r, err := resources.Init[resources.Image](expectedImageName, path)
	require.NoError(t, err, "should not error")

	signature1, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, signature1, "empty image signature")
	signature1b, err := calcImageSignature(r)
	assert.Equal(t, signature1, signature1b, "signature should stay the same")

	// Change Buildfile shoud change signature
	err = os.WriteFile(r.BuildFile, []byte("foo"), 0644)
	require.NoError(t, err, "should not error")

	signature2, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, signature2, "empty image signature")
	assert.NotEqual(t, signature1, signature2, "two signatures should differ changing buildfile")
	signature2b, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.Equal(t, signature2, signature2b, "signature should stay the same")

	// Add empty source dir shoud change signature
	subSrcDir := filepath.Join(r.AbsSourceDir(), "subSrcDir")
	err = os.MkdirAll(subSrcDir, 0755)
	require.NoError(t, err, "should not error")

	signature3a, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, signature3a, "empty image signature")
	assert.NotEqual(t, signature1, signature3a, "two signatures should differ adding empty source dir")
	assert.NotEqual(t, signature2, signature3a, "two signatures should differ adding empty source dir")

	// Add source file shoud change signature
	srcFile := filepath.Join(r.AbsSourceDir(), "srcFile")
	err = os.WriteFile(srcFile, []byte("foo"), 0644)
	require.NoError(t, err, "should not error")

	signature3b, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, signature3b, "empty image signature")
	assert.NotEqual(t, signature1, signature3b, "two signatures should differ adding empty source file")
	assert.NotEqual(t, signature2, signature3b, "two signatures should differ adding empty source file")
	assert.NotEqual(t, signature3a, signature3b, "two signatures should differ adding empty source file")

	// Add source sub dir file shoud change signature
	subSrcFile := filepath.Join(subSrcDir, "srcFile")
	err = os.WriteFile(subSrcFile, []byte("foo"), 0644)
	require.NoError(t, err, "should not error")

	signature3c, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, signature3c, "empty image signature")
	assert.NotEqual(t, signature1, signature3c, "two signatures should differ adding empty source file in sub dir")
	assert.NotEqual(t, signature2, signature3c, "two signatures should differ adding empty source file in sub dir")
	assert.NotEqual(t, signature3a, signature3c, "two signatures should differ adding empty source file in sub dir")
	assert.NotEqual(t, signature3b, signature3c, "two signatures should differ adding empty source file in sub dir")

	// Change source file shoud change signature
	err = os.WriteFile(srcFile, []byte("bar"), 0644)
	require.NoError(t, err, "should not error")

	signature3, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, signature3, "empty image signature")
	assert.NotEqual(t, signature1, signature3, "two signatures should differ changing source file")
	assert.NotEqual(t, signature2, signature3, "two signatures should differ changing source file")
	assert.NotEqual(t, signature3a, signature3, "two signatures should differ changing source file")
	assert.NotEqual(t, signature3b, signature3, "two signatures should differ changing source file")
	assert.NotEqual(t, signature3c, signature3, "two signatures should differ changing source file")

	// Change non sources shoud not change signature
	notSrcFile := filepath.Join(r.Dir(), "notSrcFile")
	err = os.WriteFile(notSrcFile, []byte("foo"), 0644)
	require.NoError(t, err, "should not error")

	signature4, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, signature3, "empty image signature")
	assert.Equal(t, signature3, signature4, "two signatures should be identical adding non source file")

	// Change in test dir shoud not change signature
	var i interface{} = r
	testable, _ := i.(resources.Tester)
	testFile := filepath.Join(testable.AbsTestDir(), "testFile")
	err = os.WriteFile(testFile, []byte("foo"), 0644)
	require.NoError(t, err, "should not error")

	signature5, err := calcImageSignature(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, signature3, "empty image signature")
	assert.Equal(t, signature3, signature5, "two signatures should be identical adding test file")
}

func TestDoesImageChanged(t *testing.T) {
	path, err := test.BuildRandTempPath()
	defer os.RemoveAll(path)
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	expectedImageName := "monImage"
	r, err := resources.Init[resources.Image](expectedImageName, path)
	require.NoError(t, err, "should not error")

	err = StoreImageSignature(r)
	require.NoError(t, err, "should not error")

	test, sign1, err := DoesImageChanged(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign1, "should not be empty")
	assert.False(t, test, "image should not be changed")

	test, sign2, err := DoesImageChanged(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign2, "should not be empty")
	assert.Equal(t, sign1, sign2, "signature should not be changed")
	assert.False(t, test, "image should not be changed")

	// Add source file shoud change signature
	srcFile := filepath.Join(r.AbsSourceDir(), "srcFile")
	err = os.WriteFile(srcFile, []byte("foo"), 0644)
	require.NoError(t, err, "should not error")

	test, sign3, err := DoesImageChanged(r)
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign3, "should not be empty")
	assert.NotEqual(t, sign1, sign3, "signature should be changed")
	assert.True(t, test, "image should be changed")
}
