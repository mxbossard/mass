package project

import (
	//"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"

	"mby.fr/utils/test"
	"mby.fr/mass/internal/workspace"
)

func initTempWorkspace(t *testing.T) (path string) {
        path, _ = test.BuildRandTempPath()
        os.Chdir(path)
        err := workspace.Init(path)
        assert.NoError(t, err, "Init should not return an error")
        return
}

func assertProjectFileTree(t *testing.T, path string) {
	assert.DirExists(t, path, "project dir file should exists")
	assert.DirExists(t, path + "/test", "test dir should exists")
	assert.FileExists(t, path + "/version.txt", "version.txt file should exists")
}

func initRandProject(t *testing.T) (name, path string) {
	name = test.RandSeq(6)
	path, _ = InitProject(name)
	assertProjectFileTree(t, path)
	return
}

func TestInitProject(t *testing.T) {
	tempDir := initTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	initRandProject(t)
}

func TestListProjects(t *testing.T) {
	tempDir := initTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	projects, _ := ListProjects()
	assert.Len(t, projects, 0, "Should list no projects")

	name, path := initRandProject(t)
	projects, _ = ListProjects()
	assert.Len(t, projects, 1, "Should list one project")
	p1 := projects[0]
	assert.Equal(t, name, p1.Name, "Bad project name")
	assert.Equal(t, path, p1.Dir, "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir, "Bad project test dir")
	assert.DirExists(t, p1.TestDir, "Project test dir does not exists")
	assert.Equal(t, defaultInitialVersion, p1.Version, "Bad project version")
	assert.Len(t, p1.Images, 0, "Should have 0 images")

	initRandProject(t)
	projects, _ = ListProjects()
	assert.Len(t, projects, 2, "Should list one project")
}

func TestInitAlreadyExistingProject(t *testing.T) {
	tempDir := initTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	name, path := initRandProject(t)
	projects, _ := ListProjects()
	assert.Len(t, projects, 1, "Should list one project")
	p1 := projects[0]
	assert.Equal(t, name, p1.Name, "Bad project name")
	assert.Equal(t, path, p1.Dir, "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir, "Bad project test dir")
	assert.DirExists(t, p1.TestDir, "Project test dir does not exists")
	assert.Equal(t, defaultInitialVersion, p1.Version, "Bad project version")

	// reinit same project
	_, err := InitProject(name)
	assert.NoError(t, err, "reiniting project should not return an error")
	projects, _ = ListProjects()
	assert.Len(t, projects, 1, "Should list one project")
	p1 = projects[0]
	assert.Equal(t, name, p1.Name, "Bad project name")
	assert.Equal(t, path, p1.Dir, "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir, "Bad project test dir")
	assert.DirExists(t, p1.TestDir, "Project test dir does not exists")
	assert.Equal(t, defaultInitialVersion, p1.Version, "Bad project version")

	// editing project
	newVersion := "0.2.1"
	os.WriteFile(path + "/version.txt", []byte(newVersion), 0644)

	// reinit edited project
	_, err = InitProject(name)
	assert.NoError(t, err, "reiniting project should not return an error")
	projects, _ = ListProjects()
	assert.Len(t, projects, 1, "Should list one project")
	p1 = projects[0]
	assert.Equal(t, name, p1.Name, "Bad project name")
	assert.Equal(t, path, p1.Dir, "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir, "Bad project test dir")
	assert.DirExists(t, p1.TestDir, "Project test dir does not exists")
	assert.Equal(t, newVersion, p1.Version, "Bad project version")
}

func TestGetNotExistingProject(t *testing.T) {
	tempDir := initTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	_, ok, err := GetProject("foo")
	assert.False(t, ok, "should not return ok")
	assert.NoError(t, err, "should not return error")
}

func TestGetExistingProject(t *testing.T) {
	tempDir := initTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	initRandProject(t)
	name, path := initRandProject(t)
	initRandProject(t)

	p, ok, err := GetProject(name)
	assert.True(t, ok, "should return ok")
	assert.NoError(t, err, "should not return error")
	assert.NotNil(t, p, "should return a project")
	assert.Equal(t, name, p.Name, "bad project name")
	assert.Equal(t, path, p.Dir, "bad project path")
}
