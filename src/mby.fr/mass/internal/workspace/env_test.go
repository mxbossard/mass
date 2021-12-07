package workspace

import (
	//"fmt"
	"testing"
	"os"
	_ "github.com/stretchr/testify/assert"

	_ "mby.fr/utils/test"
)

func TestInitEnv(t *testing.T) {
	tempDir := TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	TestInitRandEnv(t)
}

//func TestReInitEnv(t *testing.T) {
//	tempDir := TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//
//	envs, _ := ListEnvs()
//	assert.Len(t, envs, 3, "Should list default envs")
//
//	name, path := TestInitRandEnv(t, tempDir)
//	envs, _ := ListEnvs()
//	assert.Len(t, envs, 4, "Bad env listing")
//	e1 := envs[0]
//	assert.Equal(t, name, e1.Name, "Bad env name")
//	assert.Equal(t, path, e1.Dir, "Bad env dir")
//
//	// reinit same project
//	_, err := InitEnv(name)
//	require.NoError(t, err, "reiniting project should not return an error")
//	envs, _ := ListEnvs()
//	assert.Len(t, envs, 4, "Bad env listing")
//	e1 := envs[0]
//	assert.Equal(t, name, e1.Name, "Bad env name")
//	assert.Equal(t, path, e1.Dir, "Bad env dir")
//}

