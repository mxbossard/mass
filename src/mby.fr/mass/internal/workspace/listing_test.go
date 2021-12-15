package workspace

import (
	//"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	//"mby.fr/utils/test"
	"mby.fr/mass/internal/project"
)

func TestListProjects(t *testing.T) {
	tempDir := TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	projects, err := ListProjects()
	require.NoError(t, err, "should not error")
	assert.Len(t, projects, 0, "Should list 0 projects")

	name, path := project.TestInitRandProject(t)
	projects, err = ListProjects()
	require.NoError(t, err, "should not error")
	assert.Len(t, projects, 1, "Should list 1 project")
	p1 := projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")

	project.TestInitRandProject(t)
	projects, err = ListProjects()
	require.NoError(t, err, "should not error")
	assert.Len(t, projects, 2, "Should list 2 project")
}

func TestListEnvs(t *testing.T) {
	tempDir := TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	envs, err := ListEnvs()
	require.NoError(t, err, "should not error")
	assert.Len(t, envs, 3, "Should list 3 envs")

	InitEnv("foo")
	envs, err = ListEnvs()
	require.NoError(t, err, "should not error")
	assert.Len(t, envs, 4, "Should list 4 envs")
}

