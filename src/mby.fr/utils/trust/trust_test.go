package trust

import (
	//"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/utils/test"
)

func assertSignatureOk(t *testing.T, actual string, err error, msg string) {
	assert.NoError(t, err, msg+" signature should not produce an error")
	assert.NotEmpty(t, actual, msg+" signature should not be empty")
}

func assertSameSignature(t *testing.T, expected, actual string, err error, msg string) {
	assertSignatureOk(t, actual, err, msg)
	assert.Equal(t, expected, actual, msg+" signature should stay the same")
}

func assertSignatureDiffer(t *testing.T, expected, actual string, err error, msg string) {
	assertSignatureOk(t, actual, err, msg)
	assert.NotEqual(t, expected, actual, msg+" signature should be different")
}

func TestSignEmptyFile(t *testing.T) {
	path, err := test.MkRandTempFile()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path)
	assert.FileExists(t, path, "Temp file should exists")

	s1, err := SignFileContent(path)
	assertSignatureOk(t, s1, err, "empty file")

	s2, err := SignFileContent(path)
	assertSameSignature(t, s1, s2, err, "empty file")
}

func TestSignNotExistingFile(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path)
	assert.NoFileExists(t, path, "Temp file should not exists")

	_, err = SignFileContent(path)
	require.Error(t, err, "should error")
}

func TestSignFile(t *testing.T) {
	path, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path)
	assert.DirExists(t, path, "Temp dir should exists")

	// Add empty file1
	file1 := filepath.Join(path, "file1")
	os.WriteFile(file1, []byte(""), 0644)
	s1a, err := SignFileContent(file1)
	assertSignatureOk(t, s1a, err, "empty file1")

	s1b, err := SignFileContent(file1)
	assertSameSignature(t, s1a, s1b, err, "empty file1")

	// Add not empty file2
	file2 := filepath.Join(path, "file2")
	os.WriteFile(file2, []byte("foo"), 0644)
	s2a, err := SignFileContent(file2)
	assertSignatureOk(t, s2a, err, "adding file2")

	s2b, err := SignFileContent(file2)
	assertSameSignature(t, s2a, s2b, err, "adding file2")

	assertSignatureDiffer(t, s1a, s2a, err, "between 2 different files")
}

func TestSignEmptyDir(t *testing.T) {
	path, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path)
	assert.DirExists(t, path, "Temp dir should exists")

	dir1 := filepath.Join(path, "dir1")
	os.Mkdir(dir1, 0755)

	s1, err := SignDirContent(dir1)
	assertSignatureOk(t, s1, err, "empty dir1")

	s2, err := SignDirContent(dir1)
	assertSameSignature(t, s1, s2, err, "empty dir1")
}

func TestSignNotExistingDir(t *testing.T) {
	path, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path)
	assert.NoFileExists(t, path, "Temp file should not exists")

	_, err = SignDirContent(path)
	require.Error(t, err, "should error")
}

func TestSignFileInsteadOfDir(t *testing.T) {
	path, err := test.MkRandTempFile()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path)
	assert.FileExists(t, path, "Temp file should exists")

	_, err = SignDirContent(path)
	require.Error(t, err, "should error")
}

func TestSignDir(t *testing.T) {
	path, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path)
	assert.DirExists(t, path, "Temp dir should exists")

	// Empty dir1
	dir1 := filepath.Join(path, "dir1")
	os.Mkdir(dir1, 0755)

	s1, err := SignDirContent(dir1)
	assertSignatureOk(t, s1, err, "empty dir1")

	s2, err := SignDirContent(dir1)
	assertSameSignature(t, s1, s2, err, "empty dir1")

	// Add empty file1
	file1 := filepath.Join(dir1, "file1")
	os.WriteFile(file1, []byte(""), 0644)
	s3, err := SignDirContent(dir1)
	assertSignatureDiffer(t, s1, s3, err, "adding empty file1")

	s4, err := SignDirContent(dir1)
	assertSameSignature(t, s3, s4, err, "adding empty file1")

	// Add not empty file2
	file2 := filepath.Join(dir1, "file2")
	os.WriteFile(file2, []byte("foo"), 0644)
	s5, err := SignDirContent(dir1)
	assertSignatureDiffer(t, s4, s5, err, "adding file2")

	s6, err := SignDirContent(dir1)
	assertSameSignature(t, s5, s6, err, "adding file2")

	// Replace file2 content
	os.WriteFile(file2, []byte("bar"), 0644)
	s7, err := SignDirContent(dir1)
	assertSignatureDiffer(t, s6, s7, err, "modifying file2")

	s8, err := SignDirContent(dir1)
	assertSameSignature(t, s7, s8, err, "modifying file2")

	// Add sub dir2
	dir2 := filepath.Join(dir1, "dir2")
	os.Mkdir(dir2, 0755)
	s9, err := SignDirContent(dir1)
	//assertSignatureDiffer(t, s8, s9, err, "adding sub dir2")
	assertSameSignature(t, s8, s9, err, "adding sub dir2")

	s10, err := SignDirContent(dir1)
	assertSameSignature(t, s7, s10, err, "adding sub dir2")

	// Add not empty file3 in sub dir2
	file3 := filepath.Join(dir2, "file3")
	os.WriteFile(file3, []byte("baz"), 0644)
	s11, err := SignDirContent(dir1)
	assertSignatureDiffer(t, s10, s11, err, "adding filie3")

	s12, err := SignDirContent(dir1)
	assertSameSignature(t, s11, s12, err, "adding file3")
}

func TestSignNotExistsContents(t *testing.T) {
	path1, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path1)

	_, err = SignFsContents([]string{path1})
	require.Error(t, err, "should error")
}

func TestSignEmptyFsContents(t *testing.T) {
	path1, err := test.MkRandTempFile()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path1)
	assert.FileExists(t, path1, "Temp file should exists")

	s1, err := SignFsContents([]string{path1})
	assertSignatureOk(t, s1, err, "empty file1")

	s2, err := SignFsContents([]string{path1})
	assertSameSignature(t, s1, s2, err, "empty file1")

	path2, err := test.MkRandTempFile()
	require.NoError(t, err, "should not error")
	defer os.RemoveAll(path2)
	assert.FileExists(t, path2, "Temp file should exists")

	s3, err := SignFsContents([]string{path2})
	assertSignatureOk(t, s3, err, "empty file2")

	s4, err := SignFsContents([]string{path2})
	assertSameSignature(t, s3, s4, err, "empty file2")

	s5, err := SignFsContents([]string{path1, path2})
	assertSignatureOk(t, s5, err, "empty file1 and file2")

	s6, err := SignFsContents([]string{path1, path2})
	assertSameSignature(t, s5, s6, err, "empty file1 and file2")

	assertSignatureDiffer(t, s1, s3, err, "file1 signature should differ from file2")
	assertSignatureDiffer(t, s1, s5, err, "file1 signature should differ from file1+file2")
	assertSignatureDiffer(t, s3, s5, err, "file2 signature should differ from file1+file2")
}

func TestSignObject(t *testing.T) {
	sign1, err := SignObject("")
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign1, "should not be empty")

	sign2, err := SignObject("foo")
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign2, "should not be empty")

	sign3, err := SignObject("foo")
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign2, "should not be empty")
	assert.Equal(t, sign2, sign3, "should be same signature")

	sign4, err := SignObject([]string{"foo"})
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign2, "should not be empty")

	sign5, err := SignObject([]string{"foo"})
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign2, "should not be empty")
	assert.Equal(t, sign4, sign5, "should be same signature")

	sign6, err := SignObject([]string{"foo", "bar"})
	require.NoError(t, err, "should not error")
	assert.NotEmpty(t, sign2, "should not be empty")
	assert.NotEqual(t, sign5, sign6, "should be different signature")
}
