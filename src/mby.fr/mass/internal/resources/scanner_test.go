package resources

import (
	"testing"
	"os"
	//"path/filepath"

	"github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"

	"mby.fr/utils/test"
)

func initRandResource(t *testing.T, parentPath string, kind Kind) (path string) {
	resDir, err := test.MkRandSubDir(parentPath)
	require.NoError(t, err, "should not error")
	res, err := Init(resDir, kind)
	require.NoError(t, err, "should not error")
	path = res.Dir()
	return
}

func TestScanBlankPath(t *testing.T) {
	path := ""
	_, err := ScanProjects(path)
	assert.Error(t, err, "should error")

	_, err = ScanImages(path)
	assert.Error(t, err, "should error")

	_, err = ScanEnvs(path)
	assert.Error(t, err, "should error")
}

func TestScanNotExistingDir(t *testing.T) {
	path := "notExistingDirZzz"
	_, err := ScanProjects(path)
	assert.Error(t, err, "should error")

	_, err = ScanImages(path)
	assert.Error(t, err, "should error")

	_, err = ScanEnvs(path)
	assert.Error(t, err, "should error")
}

func TestScanProjects(t *testing.T) {
	parentPath, err := test.MkRandTempDir()
        defer os.RemoveAll(parentPath)

	res, err := ScanProjects(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 0, "bad resource count")

	r1Path := initRandResource(t, parentPath, ProjectKind)
	_ = r1Path
	r2Path := initRandResource(t, parentPath, ProjectKind)
	_ = r2Path
	r3Path := initRandResource(t, parentPath, ProjectKind)
	_ = r3Path
	r4Path := initRandResource(t, parentPath, EnvKind)
	_ = r4Path
	r5Path := initRandResource(t, parentPath, ImageKind)
	_ = r5Path

	// Empty dirs
	test.MkRandSubDir(parentPath)
	test.MkRandSubDir(parentPath)

	res, err = ScanProjects(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")

	res2, err := ScanImages(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res2, 1, "bad resource count")

	res3, err := ScanEnvs(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res3, 1, "bad resource count")
}

func TestScanEnvs(t *testing.T) {
	parentPath, err := test.MkRandTempDir()
        defer os.RemoveAll(parentPath)

	res, err := ScanEnvs(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 0, "bad resource count")

	r1Path := initRandResource(t, parentPath, EnvKind)
	_ = r1Path
	r2Path := initRandResource(t, parentPath, EnvKind)
	_ = r2Path
	r3Path := initRandResource(t, parentPath, EnvKind)
	_ = r3Path
	r4Path := initRandResource(t, parentPath, ProjectKind)
	_ = r4Path
	r5Path := initRandResource(t, parentPath, ImageKind)
	_ = r5Path

	// Empty dirs
	test.MkRandSubDir(parentPath)
	test.MkRandSubDir(parentPath)

	res, err = ScanEnvs(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")

	res2, err := ScanProjects(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res2, 1, "bad resource count")

	res3, err := ScanImages(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res3, 1, "bad resource count")
}

func TestScanImages(t *testing.T) {
	parentPath, err := test.MkRandTempDir()
        defer os.RemoveAll(parentPath)

	res, err := ScanImages(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 0, "bad resource count")

	r1Path := initRandResource(t, parentPath, ImageKind)
	_ = r1Path
	r2Path := initRandResource(t, parentPath, ImageKind)
	_ = r2Path
	r3Path := initRandResource(t, parentPath, ImageKind)
	_ = r3Path
	r4Path := initRandResource(t, parentPath, EnvKind)
	_ = r4Path
	r5Path := initRandResource(t, parentPath, ProjectKind)
	_ = r5Path

	// Empty dirs
	test.MkRandSubDir(parentPath)
	test.MkRandSubDir(parentPath)

	res, err = ScanImages(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")

	res2, err := ScanProjects(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res2, 1, "bad resource count")

	res3, err := ScanEnvs(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res3, 1, "bad resource count")
}
