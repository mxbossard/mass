package workspace

import (
	//"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	//"mby.fr/utils/test"
	"mby.fr/mass/internal/commontest"
	"mby.fr/mass/internal/resources"
)

func TestInitProject(t *testing.T) {
	tempDir := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	_, path := commontest.InitRandProject(t, tempDir)
	commontest.AssertProjectFileTree(t, path)
}

func TestReInitProject(t *testing.T) {
	tempDir := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(tempDir)

	name, path := commontest.InitRandProject(t, tempDir)
	commontest.AssertProjectFileTree(t, path)
	projects, err := resources.ListProjects()
	require.NoError(t, err, "should not error")
	require.NotEmpty(t, projects, "should not be empty")
	assert.Len(t, projects, 1, "Should list one project")
	p1 := projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")
	assert.DirExists(t, p1.Dir(), "Project dir does not exists")

	// reinit same project
	path, err = InitProject(name)
	commontest.AssertProjectFileTree(t, path)
	assert.NoError(t, err, "reiniting project should not return an error")
	projects, err = resources.ListProjects()
	require.NoError(t, err, "should not error")
	assert.Len(t, projects, 1, "Should list one project")
	p1 = projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")
	assert.DirExists(t, p1.Dir(), "Project dir does not exists")

	// editing project
	//newVersion := "0.2.1"
	//p1.Version = newVersion
	err = resources.Write(p1)
	require.NoError(t, err, "should not error")

	// reinit edited project
	path, err = InitProject(name)
	commontest.AssertProjectFileTree(t, path)
	assert.NoError(t, err, "reiniting project should not return an error")
	projects, err = resources.ListProjects()
	require.NoError(t, err, "should not error")
	assert.Len(t, projects, 1, "Should list one project")
	p1 = projects[0]
	assert.Equal(t, name, p1.Name(), "Bad project name")
	assert.Equal(t, path, p1.Dir(), "Bad project dir")
	assert.DirExists(t, p1.Dir(), "Project dir does not exists")
}

//func TestInitImage(t *testing.T) {
//	tempDir := TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//	name, path := TestInitRandProject(t)
//
//	p, _, _ := GetProject(name)
//	assert.Len(t, p.Images, 0, "No Image should be listed")
//
//	// Init new image
//	imageName, imagePath := TestInitRandImage(t, p.Dir())
//	assertImageFileTree(t, imagePath)
//
//	p, _, _ = GetProject(name)
//	assert.Len(t, p.Images, 1, "No Image should be listed")
//	images, err := p.Images()
//	image := images[0]
//	require.NoError(t, err, "should not error")
//	assert.Equal(t, imageName, image.Name(), "bad image name")
//	assert.Equal(t, path + "/" + imageName, image.Dir(), "bad image dir")
//	assert.DirExists(t, image.Dir(), "image dir does not exists")
//	assert.Equal(t, imagePath + "/src", image.AbsSourceDir(), "bad image source dir")
//	assert.DirExists(t, image.AbsSourceDir(), "image source dir does not exists")
//	assert.Equal(t, imagePath + "/test", image.AbsTestDir(), "bad image test dir")
//	assert.DirExists(t, image.AbsTestDir(), "image test dir does not exists")
//	assert.Equal(t, imagePath + "/Dockerfile", image.Buildfile, "bad image build file")
//	assert.FileExists(t, image.Buildfile, "image buildfile does not exists")
//	assert.Equal(t, resources.DefaultInitialVersion, image.Version, "bad image version")
//}
//
//func TestInitImages(t *testing.T) {
//	tempDir := TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//	name1, _ := TestInitRandProject(t)
//	name2, _ := TestInitRandProject(t)
//
//	p1, _, _ := GetProject(name1)
//	assert.Len(t, p1.Images, 0, "No Image should be listed")
//
//	p2, _, _ := GetProject(name2)
//	assert.Len(t, p2.Images, 0, "No Image should be listed")
//
//	// Init new image in p1
//	// Init new image
//	imagePath, _ := TestInitRandImage(t, p1.Dir())
//	assertImageFileTree(t, imagePath)
//
//	// Init new image in p1
//	imagePath, _ = TestInitRandImage(t, p1.Dir())
//	assertImageFileTree(t, imagePath)
//
//	// Init new image in p1
//	imagePath, _ = TestInitRandImage(t, p1.Dir())
//	assertImageFileTree(t, imagePath)
//
//	p1, _, _ = GetProject(name1)
//	assert.Len(t, p1.Images, 3, "Bad image count listed")
//
//	// Init new image in p2
//	imagePath, _ = TestInitRandImage(t, p2.Dir())
//	assertImageFileTree(t, imagePath)
//
//	p2, _, _ = GetProject(name2)
//	assert.Len(t, p2.Images, 1, "Bad image count listed")
//
//}
//
//func TestReInitImage(t *testing.T) {
//	tempDir := TestInitTempWorkspace(t)
//	defer os.RemoveAll(tempDir)
//	name1, _ := TestInitRandProject(t)
//	p1, _, _ := GetProject(name1)
//	assert.Len(t, p1.Images, 0, "No Image should be listed")
//
//	// Init new image in p1
//	image1 := test.RandSeq(6)
//	image1Path, err := InitImage(p1.Dir(), image1)
//	assert.NoError(t, err, "should not produce an error")
//	assertImageFileTree(t, image1Path)
//
//	// editing image
//	newVersion := "0.3.2"
//	os.WriteFile(image1Path + "/version.txt", []byte(newVersion), 0644)
//
//	// re init image
//	image1Path, err = InitImage(p1.Dir(), image1)
//	assert.NoError(t, err, "reinit should not produce an error")
//	assertImageFileTree(t, image1Path)
//
//	p1, _, _ = GetProject(name1)
//	images, err := p1.Images()
//	assert.NoError(t, err, "should not error")
//	assert.Len(t, images, 1, "project should contain only 1 image")
//	i1 := images[0]
//	assert.Equal(t, newVersion, i1.Version, "reinit image should not change version")
//}
