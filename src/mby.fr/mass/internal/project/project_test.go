package project

import (
	//"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/utils/test"
	"mby.fr/mass/internal/workspace"
	"mby.fr/mass/internal/resources"
)

func TestInitProject(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	_, path := TestInitRandProject(t)
	assertProjectFileTree(t, path)
}

//func TestListProjects(t *testing.T) {
//	tempDir := workspace.TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//
//	projects, _ := ListProjects()
//	assert.Len(t, projects, 0, "Should list no projects")
//
//	name, path := TestInitRandProject(t)
//	assertProjectFileTree(t, path)
//	projects, _ = ListProjects()
//	assert.Len(t, projects, 1, "Should list one project")
//	p1 := projects[0]
//	assert.Equal(t, name, p1.Name, "Bad project name")
//	assert.Equal(t, path, p1.Dir, "Bad project dir")
//	assert.Equal(t, path + "/test", p1.TestDir, "Bad project test dir")
//	assert.DirExists(t, p1.TestDir, "Project test dir does not exists")
//	assert.Equal(t, DefaultInitialVersion, p1.Version, "Bad project version")
//	assert.Len(t, p1.Images, 0, "Should have 0 images")
//
//	_, path = TestInitRandProject(t)
//	assertProjectFileTree(t, path)
//	projects, _ = ListProjects()
//	assert.Len(t, projects, 2, "Should list one project")
//}

func TestReInitProject(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	name, path := TestInitRandProject(t)
	assertProjectFileTree(t, path)
	projects, _ := workspace.ListProjects()
	assert.Len(t, projects, 1, "Should list one project")
	p1 := projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir(), "Bad project test dir")
	assert.DirExists(t, p1.TestDir(), "Project test dir does not exists")
	//assert.Equal(t, resources.DefaultInitialVersion, p1.Version, "Bad project version")

	// reinit same project
	path, err := InitProject(name)
	assertProjectFileTree(t, path)
	assert.NoError(t, err, "reiniting project should not return an error")
	projects, _ = workspace.ListProjects()
	assert.Len(t, projects, 1, "Should list one project")
	p1 = projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir(), "Bad project test dir")
	assert.DirExists(t, p1.TestDir(), "Project test dir does not exists")
	//assert.Equal(t, DefaultInitialVersion, p1.Version, "Bad project version")

	// editing project
	newVersion := "0.2.1"
	os.WriteFile(path + "/version.txt", []byte(newVersion), 0644)

	// reinit edited project
	path, err = InitProject(name)
	assertProjectFileTree(t, path)
	assert.NoError(t, err, "reiniting project should not return an error")
	projects, _ = workspace.ListProjects()
	assert.Len(t, projects, 1, "Should list one project")
	p1 = projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")
	assert.Equal(t, path + "/test", p1.TestDir(), "Bad project test dir")
	assert.DirExists(t, p1.TestDir(), "Project test dir does not exists")
	//assert.Equal(t, newVersion, p1.Version, "Bad project version")
}

//func TestGetNotExistingProject(t *testing.T) {
//	tempDir := workspace.TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//
//	_, ok, err := GetProject("foo")
//	assert.False(t, ok, "should not return ok")
//	assert.NoError(t, err, "should not return error")
//}
//
//func TestGetExistingProject(t *testing.T) {
//	tempDir := workspace.TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//
//	path, err := TestInitRandProject(t)
//	assertProjectFileTree(t, path)
//
//	name, path := TestInitRandProject(t)
//	path, err = TestInitRandProject(t)
//	assertProjectFileTree(t, path)
//
//	p, ok, err := GetProject(name)
//	assert.True(t, ok, "should return ok")
//	assert.NoError(t, err, "should not return error")
//	assert.NotNil(t, p, "should return a project")
//	assert.Equal(t, name, p.Name, "bad project name")
//	assert.Equal(t, path, p.Dir, "bad project path")
//}

func TestInitImage(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)
	name, path := TestInitRandProject(t)

	p, _, _ := workspace.GetProject(name)
	assert.Len(t, p.Images, 0, "No Image should be listed")

	// Init new image
	imageName, imagePath := TestInitRandImage(t, p.Dir())
	assertImageFileTree(t, imagePath)

	p, _, _ = workspace.GetProject(name)
	assert.Len(t, p.Images, 1, "No Image should be listed")
	images, err := p.Images()
	image := images[0]
	require.NoError(t, err, "should not error")
	assert.Equal(t, imageName, image.Name(), "bad image name")
	assert.Equal(t, path + "/" + imageName, image.Dir(), "bad image dir")
	assert.DirExists(t, image.Dir(), "image dir does not exists")
	assert.Equal(t, imagePath + "/src", image.SourceDir(), "bad image source dir")
	assert.DirExists(t, image.SourceDir(), "image source dir does not exists")
	assert.Equal(t, imagePath + "/test", image.TestDir(), "bad image test dir")
	assert.DirExists(t, image.TestDir(), "image test dir does not exists")
	assert.Equal(t, imagePath + "/Dockerfile", image.Buildfile, "bad image build file")
	assert.FileExists(t, image.Buildfile, "image buildfile does not exists")
	assert.Equal(t, resources.DefaultInitialVersion, image.Version, "bad image version")
}

func TestInitImages(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)
	name1, _ := TestInitRandProject(t)
	name2, _ := TestInitRandProject(t)

	p1, _, _ := workspace.GetProject(name1)
	assert.Len(t, p1.Images, 0, "No Image should be listed")

	p2, _, _ := workspace.GetProject(name2)
	assert.Len(t, p2.Images, 0, "No Image should be listed")

	// Init new image in p1
	// Init new image
	imagePath, _ := TestInitRandImage(t, p1.Dir())
	assertImageFileTree(t, imagePath)

	// Init new image in p1
	imagePath, _ = TestInitRandImage(t, p1.Dir())
	assertImageFileTree(t, imagePath)

	// Init new image in p1
	imagePath, _ = TestInitRandImage(t, p1.Dir())
	assertImageFileTree(t, imagePath)

	p1, _, _ = workspace.GetProject(name1)
	assert.Len(t, p1.Images, 3, "Bad image count listed")

	// Init new image in p2
	imagePath, _ = TestInitRandImage(t, p2.Dir())
	assertImageFileTree(t, imagePath)

	p2, _, _ = workspace.GetProject(name2)
	assert.Len(t, p2.Images, 1, "Bad image count listed")

}

func TestReInitImage(t *testing.T) {
	tempDir := workspace.TestInitTempWorkspace(t)
	defer os.RemoveAll(tempDir)
	name1, _ := TestInitRandProject(t)
	p1, _, _ := workspace.GetProject(name1)
	assert.Len(t, p1.Images, 0, "No Image should be listed")

	// Init new image in p1
	image1 := test.RandSeq(6)
	image1Path, err := InitImage(p1.Dir(), image1)
	assert.NoError(t, err, "should not produce an error")
	assertImageFileTree(t, image1Path)

	// editing image
	newVersion := "0.3.2"
	os.WriteFile(image1Path + "/version.txt", []byte(newVersion), 0644)

	// re init image
	image1Path, err = InitImage(p1.Dir(), image1)
	assert.NoError(t, err, "reinit should not produce an error")
	assertImageFileTree(t, image1Path)

	p1, _, _ = workspace.GetProject(name1)
	images, err := p1.Images()
	assert.NoError(t, err, "should not error")
	assert.Len(t, images, 1, "project should contain only 1 image")
	i1 := images[0]
	assert.Equal(t, newVersion, i1.Version, "reinit image should not change version")
}

