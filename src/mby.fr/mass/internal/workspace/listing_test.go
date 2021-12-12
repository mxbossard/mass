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

	projects, _ := ListProjects()
	assert.Len(t, projects, 0, "Should list 0 projects")

	name, path := project.TestInitRandProject(t)
	projects, _ = ListProjects()
	assert.Len(t, projects, 1, "Should list 1 project")
	p1 := projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir(), "Bad project test dir")
	assert.DirExists(t, p1.TestDir(), "Project test dir does not exists")
	//assert.Equal(t, DefaultInitialVersion, p1.Version(), "Bad project version")
	images, err := p1.Images()
	require.NoError(t, err, "should not error")
	assert.Len(t, images, 0, "Should have 0 images")

	project.TestInitRandProject(t)
	projects, _ = ListProjects()
	assert.Len(t, projects, 2, "Should list 2 project")
}

//func TestGetNotExistingProject(t *testing.T) {
//	tempDir := TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//
//	_, ok, err := GetProject("foo")
//	assert.False(t, ok, "should not return ok")
//	assert.NoError(t, err, "should not return error")
//}
//
//func TestGetExistingProject(t *testing.T) {
//	tempDir := TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//
//	TestInitRandProject(t)
//	name, path := TestInitRandProject(t)
//	TestInitRandProject(t)
//
//	p, ok, err := GetProject(name)
//	assert.True(t, ok, "should return ok")
//	assert.NoError(t, err, "should not return error")
//	assert.NotNil(t, p, "should return a project")
//	assert.Equal(t, name, p.Name, "bad project name")
//	assert.Equal(t, path, p.Dir, "bad project path")
//}

