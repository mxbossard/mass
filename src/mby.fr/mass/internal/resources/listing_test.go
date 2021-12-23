package resources

import (
	//"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/commontest"
)

func TestListProjects(t *testing.T) {
	tempDir := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	projects, err := ListProjects()
	require.NoError(t, err, "should not error")
	assert.Len(t, projects, 0, "Should list 0 projects")

	name, path := commontest.InitRandProject(t, tempDir)
	projects, err = ListProjects()
	require.NoError(t, err, "should not error")
	assert.Len(t, projects, 1, "Should list 1 project")
	p1 := projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")

	commontest.InitRandProject(t, tempDir)
	projects, err = ListProjects()
	require.NoError(t, err, "should not error")
	assert.Len(t, projects, 2, "Should list 2 project")
}

func TestListEnvs(t *testing.T) {
	tempDir := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	envs, err := ListEnvs()
	require.NoError(t, err, "should not error")
	assert.Len(t, envs, 3, "Should list 3 envs")

	commontest.InitRandEnv(t, tempDir)
	envs, err = ListEnvs()
	require.NoError(t, err, "should not error")
	assert.Len(t, envs, 4, "Should list 4 envs")
}

