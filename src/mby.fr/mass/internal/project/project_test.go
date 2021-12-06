package project

import (
	//"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"

	"mby.fr/utils/test"
	"mby.fr/mass/internal/workspace"
)

func TestInitProject(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	TestInitRandProject(t)
}

func TestListProjects(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	projects, _ := ListProjects()
	assert.Len(t, projects, 0, "Should list no projects")

	name, path := TestInitRandProject(t)
	projects, _ = ListProjects()
	assert.Len(t, projects, 1, "Should list one project")
	p1 := projects[0]
	assert.Equal(t, name, p1.Name, "Bad project name")
	assert.Equal(t, path, p1.Dir, "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir, "Bad project test dir")
	assert.DirExists(t, p1.TestDir, "Project test dir does not exists")
	assert.Equal(t, defaultInitialVersion, p1.Version, "Bad project version")
	assert.Len(t, p1.Images, 0, "Should have 0 images")

	TestInitRandProject(t)
	projects, _ = ListProjects()
	assert.Len(t, projects, 2, "Should list one project")
}

func TestReInitProject(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	name, path := TestInitRandProject(t)
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
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	_, ok, err := GetProject("foo")
	assert.False(t, ok, "should not return ok")
	assert.NoError(t, err, "should not return error")
}

func TestGetExistingProject(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	TestInitRandProject(t)
	name, path := TestInitRandProject(t)
	TestInitRandProject(t)

	p, ok, err := GetProject(name)
	assert.True(t, ok, "should return ok")
	assert.NoError(t, err, "should not return error")
	assert.NotNil(t, p, "should return a project")
	assert.Equal(t, name, p.Name, "bad project name")
	assert.Equal(t, path, p.Dir, "bad project path")
}

func TestInitImage(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)
	name, path := TestInitRandProject(t)

	p, _, _ := GetProject(name)
	assert.Len(t, p.Images, 0, "No Image should be listed")
	imageName := test.RandSeq(6)

	// Init new image
	imagePath, err := InitImage(p, imageName)
	assert.NoError(t, err, "should not produce an error")

	p, _, _ = GetProject(name)
	assert.Len(t, p.Images, 1, "No Image should be listed")
	image := p.Images[0]
	assert.Equal(t, imageName, image.Name, "bad image name")
	assert.Equal(t, path + "/" + imageName, image.Dir, "bad image dir")
	assert.DirExists(t, image.Dir, "image dir does not exists")
	assert.Equal(t, imagePath + "/src", image.SourceDir, "bad image source dir")
	assert.DirExists(t, image.SourceDir, "image source dir does not exists")
	assert.Equal(t, imagePath + "/test", image.TestDir, "bad image test dir")
	assert.DirExists(t, image.TestDir, "image test dir does not exists")
	assert.Equal(t, imagePath + "/Dockerfile", image.Buildfile, "bad image build file")
	assert.FileExists(t, image.Buildfile, "image buildfile does not exists")
	assert.Equal(t, defaultInitialVersion, image.Version, "bad image version")
}

func TestInitImages(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)
	name1, _ := TestInitRandProject(t)
	name2, _ := TestInitRandProject(t)

	p1, _, _ := GetProject(name1)
	assert.Len(t, p1.Images, 0, "No Image should be listed")

	p2, _, _ := GetProject(name2)
	assert.Len(t, p2.Images, 0, "No Image should be listed")

	// Init new image in p1
	image1 := test.RandSeq(6)
	_, err := InitImage(p1, image1)
	assert.NoError(t, err, "should not produce an error")

	// Init new image in p1
	image2 := test.RandSeq(6)
	_, err = InitImage(p1, image2)
	assert.NoError(t, err, "should not produce an error")

	// Init new image in p1
	image3 := test.RandSeq(6)
	_, err = InitImage(p1, image3)
	assert.NoError(t, err, "should not produce an error")

	p1, _, _ = GetProject(name1)
	assert.Len(t, p1.Images, 3, "Bad image count listed")

	// Init new image in p2
	image4 := test.RandSeq(6)
	_, err = InitImage(p2, image4)
	assert.NoError(t, err, "should not produce an error")

	p2, _, _ = GetProject(name2)
	assert.Len(t, p2.Images, 1, "Bad image count listed")

}

func TestReInitImage(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)
	name1, _ := TestInitRandProject(t)
	p1, _, _ := GetProject(name1)
	assert.Len(t, p1.Images, 0, "No Image should be listed")

	// Init new image in p1
	image1 := test.RandSeq(6)
	image1Path, err := InitImage(p1, image1)
	assert.NoError(t, err, "should not produce an error")

	// editing image
	newVersion := "0.3.2"
	os.WriteFile(image1Path + "/version.txt", []byte(newVersion), 0644)

	// re init image
	image1Path, err = InitImage(p1, image1)
	assert.NoError(t, err, "reinit should not produce an error")

	p1, _, _ = GetProject(name1)
	assert.Len(t, p1.Images, 1, "project should contain only 1 image")
	i1 := p1.Images[0]
	assert.Equal(t, newVersion, i1.Version, "reinit image should not change version")
}

