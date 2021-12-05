package project

import (
	"os"
	"io/fs"
	"path/filepath"

	"mby.fr/mass/internal/workspace"
)

const defaultSourceDir = "src"
const defaultTestDir = "test"
const defaultVersionFile = "version.txt"
const defaultInitialVersion = "0.0.1"
const defaultBuildFile = "Dockerfile"

var forbiddenNames = []string{defaultSourceDir, defaultTestDir, "config"}

type Image struct {
	Name string
	Dir string
	SourceDir string
	TestDir string
	Buildfile string
	Version string
}

type Project struct {
	Name string
	Dir string
	TestDir string
	Version string
	Images []Image
}

func InitProject(name string) (projectPath string, err error) {
	settingsService, err := workspace.GetSettingsService()
	if err != nil {
		return
	}

	// Create project dir
	projectPath, err = workspace.CreateSubDirectory(settingsService.ProjectsDir(), name)
	if err != nil {
		return
	}

	// Create test dir
	testDir := testDirpath(projectPath)
	err = workspace.CreateDirectory(testDir)
	if err != nil {
		return
	}

	// Init version file
	versionFile := versionFilepath(projectPath)
	_, err = softInitFile(versionFile, defaultInitialVersion)

	return
}

func InitImage(p Project, name string) (imagePath string, err error) {
	// Create image dir
	imagePath, err = workspace.CreateSubDirectory(p.Dir, name)
	if err != nil {
		return
	}

	// Create image source dir
	sourceDir := sourceDirpath(imagePath)
	err = workspace.CreateDirectory(sourceDir)
	if err != nil {
		return
	}

	// Create image test dir
	testDir := testDirpath(imagePath)
	err = workspace.CreateDirectory(testDir)
	if err != nil {
		return
	}

	// Init version file
	versionFile := versionFilepath(imagePath)
	_, err = softInitFile(versionFile, defaultInitialVersion)

	// Init Build file
	buildfile := buildfileFilepath(imagePath)
	_, err = softInitFile(buildfile, "")

	return
}

func ListProjects() (projects []Project, err error) {
	// List directories with a version file to build project list
        settingsService, err := workspace.GetSettingsService()
        if err != nil {
                return
        }

	projectsDir := settingsService.ProjectsDir()

	// Look for a version file then stop walking the branch
	projectCollector := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == defaultVersionFile {
			// Found a version file
			parentDir := filepath.Dir(path)

			// => Found a project
			p, err := buildProject(parentDir)
			if err != nil {
				return err
			}
			projects = append(projects, p)

			// Stop walking the branch
			return fs.SkipDir
		}
		return nil
	       }

	filepath.WalkDir(projectsDir, projectCollector)

	return
}

func GetProject(name string) (p Project, ok bool, err error) {
	projects, err := ListProjects()
	for _, p = range projects {
		if p.Name == name {
			ok = true
			return
		}
	}
	return
}

func listImages(p Project) (images []Image, err error) {
        // Look for a Build file then stop walking the branch
        imageCollector := func(path string, d fs.DirEntry, err error) error {
                if err != nil {
                        return err
                }
                if d.Name() == defaultBuildFile {
                        // Found a version file
                        parentDir := filepath.Dir(path)

                        // => Found an image
                        i, err := buildImage(parentDir)
                        if err != nil {
                                return err
                        }
                        images = append(images, i)

                        // Stop walking the branch
                        return fs.SkipDir
                }
                return nil
               }

        filepath.WalkDir(p.Dir, imageCollector)

        return
}

func softInitFile(filepath, content string) (path string, err error) {
	_, err = os.Stat(filepath); 
	if os.IsNotExist(err) {
		// Do not overwrite file if it already exists
		err = os.WriteFile(filepath, []byte(content), 0644)
	}
	return
}

func sourceDirpath(parentPath string) string {
	return filepath.Join(parentPath, defaultSourceDir)
}

func testDirpath(parentPath string) string {
	return filepath.Join(parentPath, defaultTestDir)
}

func versionFilepath(parentPath string) string {
	return filepath.Join(parentPath, defaultVersionFile)
}

func buildfileFilepath(parentPath string) string {
	return filepath.Join(parentPath, defaultBuildFile)
}

func buildProject(path string) (p Project, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	name := filepath.Base(path)

	testDir := testDirpath(path)
	_, err = os.Stat(testDir); 
	if os.IsNotExist(err) {
		testDir = ""
	} else if err != nil {
		return
	}

	versionFile := versionFilepath(path)
	content, err := os.ReadFile(versionFile); 
	version := string(content)
	if os.IsNotExist(err) {
		version = ""
	} else if err != nil {
		return
	}

	p = Project{Name: name, Dir: path, TestDir: testDir, Version: version}
	images, err := listImages(p)
	p.Images = images
	return
}

func buildImage(path string) (i Image, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	name := filepath.Base(path)

	sourceDir := sourceDirpath(path)
	_, err = os.Stat(sourceDir); 
	if os.IsNotExist(err) {
		sourceDir = ""
	} else if err != nil {
		return
	}

	testDir := testDirpath(path)
	_, err = os.Stat(testDir); 
	if os.IsNotExist(err) {
		testDir = ""
	} else if err != nil {
		return
	}

	buildfile := buildfileFilepath(path)
	_, err = os.Stat(buildfile); 
	if os.IsNotExist(err) {
		buildfile = ""
	} else if err != nil {
		return
	}

	versionFile := versionFilepath(path)
	content, err := os.ReadFile(versionFile); 
	version := string(content)
	if os.IsNotExist(err) {
		version = ""
	} else if err != nil {
		return
	}

	i = Image{Name: name, Dir: path, SourceDir: sourceDir, TestDir: testDir, Buildfile: buildfile, Version: version}
	return
}
