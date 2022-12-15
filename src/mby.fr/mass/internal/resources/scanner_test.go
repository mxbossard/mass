package resources

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/utils/test"
)

func initRandResource(t *testing.T, rootPath, name string, kind Kind) (path string) {
	resDir, err := test.MkRandSubDir(rootPath)
	require.NoError(t, err, "should not error")
	res, err := InitResourcer(kind, resDir, name)
	require.NoError(t, err, "should not error")
	path = res.Dir()
	return
}

func TestPathDepth(t *testing.T) {
	assert.Equal(t, 0, pathDepth("."))
	assert.Equal(t, 0, pathDepth("./"))
	assert.Equal(t, 0, pathDepth(""))
	assert.Equal(t, 0, pathDepth("foo"))
	assert.Equal(t, 0, pathDepth("foo/"))
	assert.Equal(t, 1, pathDepth("foo/bar"))
	assert.Equal(t, 1, pathDepth("foo/bar/"))
	assert.Equal(t, 2, pathDepth("foo/bar/baz"))
	assert.Equal(t, 2, pathDepth("foo/bar/baz/"))
}

func TestScanBlankPath(t *testing.T) {
	path := ""
	projects, err := Scan[Project](path)
	require.NoError(t, err, "should not error")
	assert.Empty(t, projects, "should be empty")

	images, err := Scan[Image](path)
	require.NoError(t, err, "should not error")
	assert.Empty(t, images, "should be empty")

	envs, err := Scan[Env](path)
	require.NoError(t, err, "should not error")
	assert.Empty(t, envs, "should be empty")
}

func TestScanNotExistingDir(t *testing.T) {
	path := "notExistingDirZzz"
	projects, err := Scan[Project](path)
	assert.NoError(t, err, "should error")
	assert.Empty(t, projects, "should be empty")

	images, err := Scan[Image](path)
	assert.NoError(t, err, "should error")
	assert.Empty(t, images, "should be empty")

	envs, err := Scan[Env](path)
	assert.NoError(t, err, "should error")
	assert.Empty(t, envs, "should be empty")
}

func TestScanProjects(t *testing.T) {
	parentPath, err := test.MkRandTempDir()
	defer os.RemoveAll(parentPath)

	res, err := Scan[Project](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 0, "bad resource count")

	r1Path := initRandResource(t, parentPath, "p1", ProjectKind)
	_ = r1Path
	r2Path := initRandResource(t, parentPath, "p2", ProjectKind)
	_ = r2Path
	r3Path := initRandResource(t, parentPath, "p3", ProjectKind)
	_ = r3Path
	r4Path := initRandResource(t, parentPath, "e1", EnvKind)
	_ = r4Path
	r5Path := initRandResource(t, parentPath, "i1", ImageKind)
	_ = r5Path

	// Empty dirs
	test.MkRandSubDir(parentPath)
	test.MkRandSubDir(parentPath)

	res, err = Scan[Project](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")

	res2, err := Scan[Image](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res2, 1, "bad resource count")

	res3, err := Scan[Env](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res3, 1, "bad resource count")

	parentDepth := 0
	res4, err := ScanMaxDepth[Project](parentPath, parentDepth+0)
	require.NoError(t, err, "should not error")
	assert.Len(t, res4, 0, "bad resource count")

	res5, err := ScanMaxDepth[Project](parentPath, parentDepth+1)
	require.NoError(t, err, "should not error")
	assert.Len(t, res5, 3, "bad resource count")

	res6, err := ScanMaxDepth[Project](parentPath, parentDepth+2)
	require.NoError(t, err, "should not error")
	assert.Len(t, res6, 3, "bad resource count")
}

func TestScanEnvs(t *testing.T) {
	parentPath, err := test.MkRandTempDir()
	defer os.RemoveAll(parentPath)

	res, err := Scan[Env](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 0, "bad resource count")

	r1Path := initRandResource(t, parentPath, "e1", EnvKind)
	_ = r1Path
	r2Path := initRandResource(t, parentPath, "e2", EnvKind)
	_ = r2Path
	r3Path := initRandResource(t, parentPath, "e3", EnvKind)
	_ = r3Path
	r4Path := initRandResource(t, parentPath, "p1", ProjectKind)
	_ = r4Path
	r5Path := initRandResource(t, parentPath, "i1", ImageKind)
	_ = r5Path

	// Empty dirs
	test.MkRandSubDir(parentPath)
	test.MkRandSubDir(parentPath)

	res, err = Scan[Env](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")

	res2, err := Scan[Project](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res2, 1, "bad resource count")

	res3, err := Scan[Image](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res3, 1, "bad resource count")

	parentDepth := 0
	res4, err := ScanMaxDepth[Env](parentPath, parentDepth+0)
	require.NoError(t, err, "should not error")
	assert.Len(t, res4, 0, "bad resource count")

	res5, err := ScanMaxDepth[Env](parentPath, parentDepth+1)
	require.NoError(t, err, "should not error")
	assert.Len(t, res5, 3, "bad resource count")

	res6, err := ScanMaxDepth[Env](parentPath, parentDepth+2)
	require.NoError(t, err, "should not error")
	assert.Len(t, res6, 3, "bad resource count")
}

func TestScanImages(t *testing.T) {
	parentPath, err := test.MkRandTempDir()
	defer os.RemoveAll(parentPath)

	res, err := Scan[Image](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 0, "bad resource count")

	r1Path := initRandResource(t, parentPath, "i1", ImageKind)
	_ = r1Path
	r2Path := initRandResource(t, parentPath, "i2", ImageKind)
	_ = r2Path
	r3Path := initRandResource(t, parentPath, "i3", ImageKind)
	_ = r3Path
	r4Path := initRandResource(t, parentPath, "e1", EnvKind)
	_ = r4Path
	r5Path := initRandResource(t, parentPath, "p1", ProjectKind)
	_ = r5Path

	// Empty dirs
	test.MkRandSubDir(parentPath)
	test.MkRandSubDir(parentPath)

	res, err = Scan[Image](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")

	res2, err := Scan[Project](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res2, 1, "bad resource count")

	res3, err := Scan[Env](parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res3, 1, "bad resource count")

	parentDepth := 0
	res4, err := ScanMaxDepth[Image](parentPath, parentDepth+0)
	require.NoError(t, err, "should not error")
	assert.Len(t, res4, 0, "bad resource count")

	res5, err := ScanMaxDepth[Image](parentPath, parentDepth+1)
	require.NoError(t, err, "should not error")
	assert.Len(t, res5, 3, "bad resource count")

	res6, err := ScanMaxDepth[Image](parentPath, parentDepth+2)
	require.NoError(t, err, "should not error")
	assert.Len(t, res6, 3, "bad resource count")
}

func TestScanResourcesFrom(t *testing.T) {
	fakeWorkspacePath := initWorkspace(t)

	expectedImagesCount := 9
	var err error
	var resources []Resourcer
	// Scan all depths from root dir
	resources, err = scanResourcesFrom(fakeWorkspacePath, EnvKind, -1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(envs), "bad resource count")

	resources, err = scanResourcesFrom(fakeWorkspacePath, ProjectKind, -1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(projects), "bad resource count")

	resources, err = scanResourcesFrom(fakeWorkspacePath, ImageKind, -1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, expectedImagesCount, "bad resource count")

	resources, err = scanResourcesFrom(fakeWorkspacePath, AllKind, -1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(envs)+len(projects)+expectedImagesCount, "bad resource count")

	// Scan 1 depth from root dir
	resources, err = scanResourcesFrom(fakeWorkspacePath, EnvKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, 0, "bad resource count")

	resources, err = scanResourcesFrom(fakeWorkspacePath, ProjectKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(projects), "bad resource count")

	resources, err = scanResourcesFrom(fakeWorkspacePath, ImageKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, 0, "bad resource count")

	resources, err = scanResourcesFrom(fakeWorkspacePath, AllKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(projects), "bad resource count")

	envDirPath := filepath.Join(fakeWorkspacePath, envDir)

	// Scan all depths from env dir
	resources, err = scanResourcesFrom(envDirPath, EnvKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(envs), "bad resource count")

	resources, err = scanResourcesFrom(envDirPath, ProjectKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, 0, "bad resource count")

	resources, err = scanResourcesFrom(envDirPath, ImageKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, 0, "bad resource count")

	resources, err = scanResourcesFrom(envDirPath, AllKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(envs), "bad resource count")

	// Scan 1 depth from env dir
	resources, err = scanResourcesFrom(envDirPath, EnvKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(envs), "bad resource count")

	resources, err = scanResourcesFrom(envDirPath, ProjectKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, 0, "bad resource count")

	resources, err = scanResourcesFrom(envDirPath, ImageKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, 0, "bad resource count")

	resources, err = scanResourcesFrom(envDirPath, AllKind, 1)
	require.NoError(t, err, "should not error")
	assert.Len(t, resources, len(envs), "bad resource count")
}
