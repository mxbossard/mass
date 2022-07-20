package resources

import (
	"os"
	"testing"
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
	projects, err := ScanProjects(path)
	require.NoError(t, err, "should not error")
	assert.Empty(t, projects, "should be empty")

	images, err := ScanImages(path)
	require.NoError(t, err, "should not error")
	assert.Empty(t, images, "should be empty")

	envs, err := ScanEnvs(path)
	require.NoError(t, err, "should not error")
	assert.Empty(t, envs, "should be empty")
}

func TestScanNotExistingDir(t *testing.T) {
	path := "notExistingDirZzz"
	projects, err := ScanProjects(path)
	assert.NoError(t, err, "should error")
	assert.Empty(t, projects, "should be empty")

	images, err := ScanImages(path)
	assert.NoError(t, err, "should error")
	assert.Empty(t, images, "should be empty")

	envs, err := ScanEnvs(path)
	assert.NoError(t, err, "should error")
	assert.Empty(t, envs, "should be empty")
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
