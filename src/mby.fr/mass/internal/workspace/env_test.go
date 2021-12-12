package workspace

import (
	//"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "mby.fr/utils/test"
)

func TestInitEnv(t *testing.T) {
	tempDir := TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	path := TestInitRandEnv(t)

	assertEnvFileTree(t, path)
}

func TestReInitEnv(t *testing.T) {
	tempDir := TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	envs, _ := ListEnvs()
	assert.Len(t, envs, 3, "Should list default envs")

	_ = TestInitRandEnv(t)
	envs, _ = ListEnvs()
	assert.Len(t, envs, 4, "Bad env listing")
	e1 := envs[0]

	// reinit env e1
	_, err := InitEnv(e1.Name())
	require.NoError(t, err, "reiniting project should not return an error")
	envs, _ = ListEnvs()
	assert.Len(t, envs, 4, "Bad env listing")
	e2 := envs[0]
	assert.Equal(t, e1.Name(), e2.Name(), "Bad env name")
	assert.Equal(t, e1.Dir(), e2.Dir(), "Bad env dir")
}

