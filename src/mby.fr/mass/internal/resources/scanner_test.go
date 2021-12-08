package resources

import (
	"testing"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"

	"mby.fr/utils/test"
)

func initRandResource(t *testing.T, parentPath, kind string) (path string) {
	resDir, err := test.MkRandSubDir(parentPath)
	require.NoError(t, err, "should not error")
	resFilepath := filepath.Join(resDir, defaultResourceFile)
	resYamlContent := "resourcekind: " + kind
	err = os.WriteFile(resFilepath, []byte(resYamlContent), 0644)
	require.NoError(t, err, "should not error")
	return
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

	res, err = ScanProjects(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")
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

	res, err = ScanEnvs(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")
}

func TestScanImages(t *testing.T) {
	parentPath, err := test.MkRandTempDir()
        defer os.RemoveAll(parentPath)

	res, err := ScanImages(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 0, "bad resource count")

	r1Path := initRandResource(t, parentPath, EnvKind)
	_ = r1Path
	r2Path := initRandResource(t, parentPath, EnvKind)
	_ = r2Path
	r3Path := initRandResource(t, parentPath, EnvKind)
	_ = r3Path

	res, err = ScanImages(parentPath)
	require.NoError(t, err, "should not error")
	assert.Len(t, res, 3, "bad resource count")
}
