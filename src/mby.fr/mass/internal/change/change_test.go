package change

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/test"
)

func TestInit(t *testing.T) {
	tempDir, err := test.MkRandTempDir()
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
	require.NoError(t, err, "should not error")

	r, err := resources.BuildImage(path)
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
	err = settings.Init(path)
	require.NoError(t, err, "should not error")
	os.Chdir(path)

	err = r.Init()
	require.NoError(t, err, "should not error")

	signature, e := calcImageSignature(r)
	require.NoError(t, e, "should not error")
	assert.NotEmpty(t, signature, "empty image signature")
}
